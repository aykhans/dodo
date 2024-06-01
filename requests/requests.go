package requests

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"

	"sync"
	"time"

	"github.com/aykhans/dodo/config"
	"github.com/aykhans/dodo/custom_errors"
	"github.com/aykhans/dodo/readers"
	"github.com/aykhans/dodo/utils"
	"github.com/jedib0t/go-pretty/v6/progress"
)

type DodoResponse struct {
	Response string
	Time     time.Duration
}

type DodoResponses []DodoResponse

type MergedDodoResponse struct {
	Response string
	Count    int
	AvgTime  time.Duration
	MinTime  time.Duration
	MaxTime  time.Duration
}

func (d DodoResponses) Len() int {
	return len(d)
}

func (d DodoResponses) MinTime() time.Duration {
	minTime := d[0].Time
	for _, response := range d {
		if response.Time < minTime {
			minTime = response.Time
		}
	}
	return minTime
}

func (d DodoResponses) MaxTime() time.Duration {
	maxTime := d[0].Time
	for _, response := range d {
		if response.Time > maxTime {
			maxTime = response.Time
		}
	}
	return maxTime
}

func (d DodoResponses) AvgTime() time.Duration {
	var sum time.Duration
	for _, response := range d {
		sum += response.Time
	}
	return sum / time.Duration(len(d))
}

func (d DodoResponses) MergeDodoResponses() []MergedDodoResponse {
	mergedResponses := make(map[string]*struct {
		count     int
		minTime   time.Duration
		maxTime   time.Duration
		totalTime time.Duration
	})
	for _, response := range d {
		if _, ok := mergedResponses[response.Response]; !ok {
			mergedResponses[response.Response] = &struct {
				count     int
				minTime   time.Duration
				maxTime   time.Duration
				totalTime time.Duration
			}{
				count:     1,
				minTime:   response.Time,
				maxTime:   response.Time,
				totalTime: response.Time,
			}
		} else {
			mergedResponses[response.Response].count++
			mergedResponses[response.Response].totalTime += response.Time
			if response.Time < mergedResponses[response.Response].minTime {
				mergedResponses[response.Response].minTime = response.Time
			}
			if response.Time > mergedResponses[response.Response].maxTime {
				mergedResponses[response.Response].maxTime = response.Time
			}

		}
	}
	var result []MergedDodoResponse
	for response, data := range mergedResponses {
		result = append(result, MergedDodoResponse{
			Response: response,
			Count:    data.count,
			AvgTime:  data.totalTime / time.Duration(data.count),
			MinTime:  data.minTime,
			MaxTime:  data.maxTime,
		})
	}
	return result
}

func Run(conf *config.DodoConfig) (DodoResponses, error) {
	params := setParams(conf.URL, conf.Params)
	headers := getHeaders(conf.Headers)

	dodosCountForRequest, dodosCountForProxies := conf.DodosCount, conf.DodosCount
	if dodosCountForRequest > conf.RequestCount {
		dodosCountForRequest = conf.RequestCount
	}
	proxiesCount := len(conf.Proxies)
	if dodosCountForProxies > proxiesCount {
		dodosCountForProxies = proxiesCount
	}
	dodosCountForProxies = min(dodosCountForProxies, config.MaxDodosCountForProxies)

	var wg sync.WaitGroup
	wg.Add(dodosCountForRequest + 1)
	var requestCountPerDodo int
	responses := make([][]DodoResponse, dodosCountForRequest)
	getClient := getClientFunc(conf.Proxies, conf.Timeout, dodosCountForProxies)

	countSlice := make([]int, dodosCountForRequest)
	go printProgress(&wg, conf.RequestCount, "Dodos Working🔥", &countSlice)

	for i := 0; i < dodosCountForRequest; i++ {
		if i+1 == dodosCountForRequest {
			requestCountPerDodo = conf.RequestCount -
				(i * conf.RequestCount / dodosCountForRequest)
		} else {
			requestCountPerDodo = ((i + 1) * conf.RequestCount / dodosCountForRequest) -
				(i * conf.RequestCount / dodosCountForRequest)
		}
		go sendRequest(
			&responses[i],
			&countSlice[i],
			requestCountPerDodo,
			conf.Method,
			params,
			conf.Body,
			headers,
			conf.Cookies,
			getClient,
			&wg,
		)
	}
	wg.Wait()
	return utils.Flatten(responses), nil
}

func sendRequest(
	responseData *[]DodoResponse,
	counter *int,
	requestCout int,
	method string,
	params string,
	body string,
	headers http.Header,
	cookies map[string]string,
	getClient func() http.Client,
	wg *sync.WaitGroup,
) {
	defer wg.Done()
	for j := 0; j < requestCout; j++ {
		func() {
			defer func() { *counter++ }()
			req, _ := http.NewRequest(
				method,
				params,
				getBodyReader(body),
			)
			req.Header = headers
			setCookies(req, cookies)
			client := getClient()
			startTime := time.Now()
			resp, err := client.Do(req)
			completedTime := time.Since(startTime)
			if err != nil {
				*responseData = append(
					*responseData,
					DodoResponse{
						Response: customerrors.RequestErrorsFormater(err),
						Time:     completedTime,
					},
				)
				return
			}
			defer resp.Body.Close()
			*responseData = append(
				*responseData,
				DodoResponse{
					Response: resp.Status,
					Time:     completedTime,
				},
			)
		}()
	}
}

func setCookies(req *http.Request, cookies map[string]string) {
	for key, value := range cookies {
		req.AddCookie(&http.Cookie{Name: key, Value: value})
	}
}

func getHeaders(headers map[string]string) http.Header {
	httpHeaders := make(http.Header, len(headers))
	httpHeaders.Set("User-Agent", config.DefaultUserAgent)
	for key, value := range headers {
		httpHeaders.Add(key, value)
	}
	return httpHeaders
}

func getBodyReader(bodyString string) io.Reader {
	if bodyString == "" {
		return http.NoBody
	}
	return strings.NewReader(bodyString)

}

func setParams(baseURL string, params map[string]string) string {
	if len(params) == 0 {
		return baseURL
	}
	urlParams := url.Values{}
	for key, value := range params {
		urlParams.Add(key, value)
	}
	baseURLWithParams := fmt.Sprintf("%s?%s", baseURL, urlParams.Encode())
	return baseURLWithParams
}

func printProgress(wg *sync.WaitGroup, total int, message string, countSlice *[]int) {
	defer wg.Done()
	pw := progress.NewWriter()
	pw.SetTrackerPosition(progress.PositionRight)
	pw.SetStyle(progress.StyleBlocks)
	pw.SetTrackerLength(40)
	pw.SetUpdateFrequency(time.Millisecond * 250)
	go pw.Render()
	dodosTracker := progress.Tracker{Message: message, Total: int64(total)}
	pw.AppendTracker(&dodosTracker)
	for {
		totalCount := 0
		for _, count := range *countSlice {
			totalCount += count
		}
		dodosTracker.SetValue(int64(totalCount))
		if totalCount == total {
			break
		}
		time.Sleep(time.Millisecond * 200)
	}
	dodosTracker.MarkAsDone()
	time.Sleep(time.Millisecond * 300)
	pw.Stop()
}

func getClientFunc(proxies []config.Proxy, timeout time.Duration, dodosCount int) func() http.Client {
	if len(proxies) > 0 {
		activeProxyClientsArray := make([][]http.Client, dodosCount)
		proxiesCount := len(proxies)
		var wg sync.WaitGroup
		wg.Add(dodosCount + 1)
		var proxiesSlice []config.Proxy

		countSlice := make([]int, dodosCount)
		go printProgress(&wg, proxiesCount, "Searching for active proxies🌐", &countSlice)

		for i := 0; i < dodosCount; i++ {
			if i+1 == dodosCount {
				proxiesSlice = proxies[i*proxiesCount/dodosCount:]
			} else {
				proxiesSlice = proxies[i*proxiesCount/dodosCount : (i+1)*proxiesCount/dodosCount]
			}
			go findActiveProxyClients(
				proxiesSlice,
				timeout,
				&activeProxyClientsArray[i],
				&countSlice[i],
				&wg,
			)
		}
		wg.Wait()

		activeProxyClients := utils.Flatten(activeProxyClientsArray)
		activeProxyClientsCount := len(activeProxyClients)
		var yesOrNoMessage string
		if activeProxyClientsCount == 0 {
			yesOrNoMessage = utils.Colored(
				utils.Colors.Red,
				"No active proxies found. Do you want to continue?",
			)
		} else {
			yesOrNoMessage = utils.Colored(
				utils.Colors.Yellow,
				fmt.Sprintf("Found %d active proxies. Do you want to continue?", activeProxyClientsCount),
			)
		}
		fmt.Println()
		proceed := readers.CLIYesOrNoReader(yesOrNoMessage)
		if !proceed {
			utils.PrintAndExit("Exiting...")
		}
		fmt.Println()
		if activeProxyClientsCount == 0 {
			return func() http.Client {
				return getNewClient(timeout)
			}
		}
		return func() http.Client {
			return getRandomClient(activeProxyClients, activeProxyClientsCount)
		}
	}
	return func() http.Client {
		return getNewClient(timeout)
	}
}

func findActiveProxyClients(
	proxies []config.Proxy,
	timeout time.Duration,
	activeProxyClients *[]http.Client,
	counter *int,
	wg *sync.WaitGroup) {
	defer wg.Done()
	for _, proxy := range proxies {
		func() {
			defer func() { *counter++ }()
			transport, err := getTransport(proxy)
			if err != nil {
				return
			}
			client := &http.Client{
				Transport: transport,
				Timeout:   timeout,
			}
			resp, err := client.Get(config.ProxyCheckURL)
			if err != nil {
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode == 200 {
				*activeProxyClients = append(
					*activeProxyClients,
					http.Client{
						Transport: transport,
						Timeout:   timeout,
					},
				)
			}
		}()
	}
}

func getTransport(proxy config.Proxy) (*http.Transport, error) {
	proxyURL, err := url.Parse(proxy.URL)
	if err != nil {
		return nil, err
	}
	if proxy.Username != "" {
		transport := &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
		return transport, nil
	}

	transport := &http.Transport{
		Proxy: http.ProxyURL(
			&url.URL{
				Scheme: proxyURL.Scheme,
				Host:   proxyURL.Host,
				User:   url.UserPassword(proxy.Username, proxy.Password),
			},
		),
	}
	return transport, nil
}

func getRandomClient(clients []http.Client, clientsCount int) http.Client {
	randomIndex := rand.Intn(clientsCount)
	return clients[randomIndex]
}

func getNewClient(timeout time.Duration) http.Client {
	return http.Client{Timeout: timeout}
}
