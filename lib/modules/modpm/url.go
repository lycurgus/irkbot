package modpm

import (
    "fmt"
    "golang.org/x/net/html"
    "net/http"
    "strings"
    urllib "net/url"
    "github.com/davidscholberg/irkbot/lib"
    "github.com/mvdan/xurls"
)

// Url attempts to fetch the title of the HTML document returned by a URL
func Url(p *lib.Privmsg) bool {
    urls := xurls.Strict.FindAllString(p.Msg, -1)

    for _, urlStr := range urls {
        url, err := urllib.Parse(urlStr)
        if err != nil {
            continue
        }
        if v, _ := validateUrl(url); !v {
            lib.Say(p, fmt.Sprintf("%s: :|", p.Event.Nick))
            break
        }
        host := url.Host
        title, err := getHtmlTitle(urlStr)
        if err != nil {
            continue
        }
        title = strings.Replace(title, "\n", "", -1)
        title = strings.Replace(title, "\r", "", -1)
        title = strings.TrimSpace(title)
        lib.Say(p, fmt.Sprintf("^ %s - [%s]", title, host))
    }

    // don't consume the message, in case there are commands in it
    return false
}

// getHtmlTitle returns the HTML title found at the given URL.
func getHtmlTitle(url string) (string, error) {
    response, err := http.Get(url)
    if err != nil {
        return "", err
    }
    defer response.Body.Close()

    // ignore response codes 400 and above
    if response.StatusCode >= 400 {
        return "", fmt.Errorf("received status %d", response.StatusCode)
    }

    doctree, err := html.Parse(response.Body)
    if err != nil {
        return "", err
    }

    title, err := searchForHtmlTitle(doctree)
    if err != nil {
        return "", err
    }
    if len(title) == 0 {
        return "", fmt.Errorf("title not found")
    }

    return title, nil
}

// validateUrl ensures that the given URL is safe to GET.
func validateUrl(url *urllib.URL) (bool, error) {
    // TODO: add CIDR matching
    if strings.HasPrefix(url.Host, "127.") ||
        strings.HasPrefix(url.Host, "192.168.") ||
        strings.HasPrefix(url.Host, "localhost") {
        return false, fmt.Errorf("host not allowed")
    }

    return true, nil
}

// searchForHtmlTitle searches the parsed html document for the title.
func searchForHtmlTitle(n *html.Node) (string, error) {
    if n.Type == html.ElementNode && n.Data == "title" {
        if n.FirstChild.Type != html.TextNode {
            err := fmt.Errorf("child of title not TextNode type")
            return "", err
        }
        return n.FirstChild.Data, nil
    }
    for c := n.FirstChild; c != nil; c = c.NextSibling {
        title, err := searchForHtmlTitle(c)
        if len(title) > 0 || err != nil {
            return title, err
        }
    }
    return "", nil
}
