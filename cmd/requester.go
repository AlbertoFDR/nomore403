package cmd

import (
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/cheynewallace/tabby"
	"github.com/fatih/color"
	"github.com/zenthangplus/goccm"
)

type Result struct {
	line          string
	statusCode    int
	contentLength int
}

func printResponse(results []Result) {
	t := tabby.New()

	var code string
	for _, result := range results {
		switch result.statusCode {
		case 200, 201, 202, 203, 204, 205, 206:
			code = color.GreenString(strconv.Itoa(result.statusCode))
		case 300, 301, 302, 303, 304, 307, 308:
			code = color.YellowString(strconv.Itoa(result.statusCode))
		case 400, 401, 402, 403, 404, 405, 406, 407, 408, 413, 429:
			code = color.RedString(strconv.Itoa(result.statusCode))
		case 500, 501, 502, 503, 504, 505, 511:
			code = color.MagentaString(strconv.Itoa(result.statusCode))
		}
		t.AddLine(code, color.BlueString(strconv.Itoa(result.contentLength)+" bytes"), result.line)
	}
	t.Print()

}

func requestMethods(uri string, headers []header, proxy *url.URL, folder string) {
	color.Cyan("\n[####] VERB TAMPERING [####]")

	var lines []string
	lines, err := parseFile(folder + "/httpmethods")
	if err != nil {
		log.Fatal(err)
	}

	w := goccm.New(max_goroutines)

	results := []Result{}

	for _, line := range lines {
		time.Sleep(time.Duration(delay) * time.Millisecond)
		w.Wait()
		go func(line string) {
			statusCode, response, err := request(line, uri, headers, proxy)
			if err != nil {
				log.Println(err)
			}

			results = append(results, Result{line, statusCode, len(response)})
			w.Done()
		}(line)
	}
	w.WaitAllDone()
	printResponse(results)
}

func requestHeaders(uri string, headers []header, proxy *url.URL, bypassIp string, folder string, method string) {
	color.Cyan("\n[####] HEADERS [####]")

	var lines []string
	lines, err := parseFile(folder + "/headers")
	if err != nil {
		log.Fatal(err)
	}

	var ips []string
	if len(bypassIp) != 0 {
		ips = []string{bypassIp}
	} else {
		ips, err = parseFile(folder + "/ips")
		if err != nil {
			log.Fatal(err)
		}
	}

	simpleheaders, err := parseFile(folder + "/simpleheaders")
	if err != nil {
		log.Fatal(err)
	}

	w := goccm.New(max_goroutines)

	results := []Result{}

	for _, ip := range ips {
		for _, line := range lines {
			time.Sleep(time.Duration(delay) * time.Millisecond)
			w.Wait()
			go func(line, ip string) {
				headers := append(headers, header{line, ip})

				statusCode, response, err := request(method, uri, headers, proxy)

				if err != nil {
					log.Println(err)
				}

				results = append(results, Result{line + ": " + ip, statusCode, len(response)})
				w.Done()
			}(line, ip)
		}
	}

	for _, simpleheader := range simpleheaders {
		time.Sleep(time.Duration(delay) * time.Millisecond)
		w.Wait()
		go func(line string) {
			x := strings.Split(line, " ")
			headers := append(headers, header{x[0], x[1]})

			statusCode, response, err := request(method, uri, headers, proxy)
			if err != nil {
				log.Println(err)
			}

			results = append(results, Result{x[0] + ": " + x[1], statusCode, len(response)})
			w.Done()
		}(simpleheader)
	}
	w.WaitAllDone()
	printResponse(results)
}

func requestEndPaths(uri string, headers []header, proxy *url.URL, folder string, method string) {
	color.Cyan("\n[####] CUSTOM PATHS [####]")

	var lines []string
	lines, err := parseFile(folder + "/endpaths")
	if err != nil {
		log.Fatal(err)
	}

	w := goccm.New(max_goroutines)

	results := []Result{}

	for _, line := range lines {
		time.Sleep(time.Duration(delay) * time.Millisecond)
		w.Wait()
		go func(line string) {
			statusCode, response, err := request(method, uri+line, headers, proxy)
			if err != nil {
				log.Println(err)
			}

			results = append(results, Result{uri + line, statusCode, len(response)})
			w.Done()
		}(line)
	}
	w.WaitAllDone()
	printResponse(results)
}

func requestMidPaths(uri string, headers []header, proxy *url.URL, folder string, method string) {
	var lines []string
	lines, err := parseFile(folder + "/midpaths")
	if err != nil {
		log.Fatal(err)
	}

	x := strings.Split(uri, "/")
	var uripath string

	if uri[len(uri)-1:] == "/" {
		uripath = x[len(x)-2]
	} else {
		uripath = x[len(x)-1]
	}

	baseuri := strings.ReplaceAll(uri, uripath, "")
	baseuri = baseuri[:len(baseuri)-1]

	w := goccm.New(max_goroutines)

	results := []Result{}

	for _, line := range lines {
		time.Sleep(time.Duration(delay) * time.Millisecond)
		w.Wait()
		go func(line string) {
			var fullpath string
			if uri[len(uri)-1:] == "/" {
				fullpath = baseuri + line + uripath + "/"
			} else {
				fullpath = baseuri + "/" + line + uripath
			}

			statusCode, response, err := request(method, fullpath, headers, proxy)
			if err != nil {
				log.Println(err)
			}

			results = append(results, Result{fullpath, statusCode, len(response)})
			w.Done()
		}(line)
	}
	w.WaitAllDone()
	printResponse(results)
}

func requestCapital(uri string, headers []header, proxy *url.URL, method string) {
	color.Cyan("\n[####] CAPITALIZATION [####]")

	x := strings.Split(uri, "/")
	var uripath string

	if uri[len(uri)-1:] == "/" {
		uripath = x[len(x)-2]
	} else {
		uripath = x[len(x)-1]
	}
	baseuri := strings.ReplaceAll(uri, uripath, "")
	baseuri = baseuri[:len(baseuri)-1]

	w := goccm.New(max_goroutines)

	results := []Result{}

	for _, z := range uripath {
		time.Sleep(time.Duration(delay) * time.Millisecond)
		w.Wait()
		go func(z string) {
			newpath := strings.ReplaceAll(uripath, string(z), strings.ToUpper(string(z)))

			var fullpath string
			if uri[len(uri)-1:] == "/" {
				fullpath = baseuri + newpath + "/"
			} else {
				fullpath = baseuri + "/" + newpath
			}

			statusCode, response, err := request(method, fullpath, headers, proxy)
			if err != nil {
				log.Println(err)
			}

			results = append(results, Result{fullpath, statusCode, len(response)})
			w.Done()
		}(string(z))
	}
	w.WaitAllDone()
	printResponse(results)
}

func requester(uri string, proxy string, userAgent string, req_headers []string, bypassIp string, folder string, method string) {
	if len(proxy) != 0 {
		if !strings.Contains(proxy, "http") {
			proxy = "http://" + proxy
		}
		color.Magenta("\n[*] USING PROXY: %s\n", proxy)
	}
	userProxy, _ := url.Parse(proxy)
	x := strings.Split(uri, "/")
	if len(x) < 4 {
		uri += "/"
	}
	if len(userAgent) == 0 {
		userAgent = "dontgo403"
	}
	if len(method) == 0 {
		method = "GET"
	}

	headers := []header{
		{"User-Agent", userAgent},
	}

	if len(req_headers[0]) != 0 {
		for _, _header := range req_headers {
			header_split := strings.Split(_header, ":")
			headers = append(headers, header{header_split[0], header_split[1]})
		}
	}

	requestMethods(uri, headers, userProxy, folder)
	requestHeaders(uri, headers, userProxy, bypassIp, folder, method)
	requestEndPaths(uri, headers, userProxy, folder, method)
	requestMidPaths(uri, headers, userProxy, folder, method)
	requestCapital(uri, headers, userProxy, method)
}
