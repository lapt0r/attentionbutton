package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/hajimehoshi/oto"

	"github.com/hajimehoshi/go-mp3"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", "attention")
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	html := `
		<html>
	<meta name="viewport" content="width=device-width, initial-scale=1.0"> 
	<head>
		<script>
			document.addEventListener("DOMContentLoaded", function(){


						const button = document.getElementById('attn-btn');

						button.addEventListener('click', async _ => {
						try {     
							const response = await fetch('http://192.168.0.22:8080/api/attention', {
							method: 'post',
							});
							console.log('Completed!', response);
						} catch(err) {
							console.error('fetch of backend service failed');
						}
						});
			});
		</script>
		<style>

		.center {
		display: flex;
		justify-content: center;
		align-items: center;
		}
		</style>
	</head>
	<div class=center><h1>
		attention
	</h1></div>

	<div class="center">
		<button id='attn-btn' style="min-width: 200px;min-height:100px;max-width:90%;max-height:50%;border-color:pink">	&#x1F9A9; attention &#x1F9A9;</button>
	</div>
	</html>
	`
	fmt.Fprint(w, html)
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	f, err := os.Open("flamingo_short.mp3")
	if err != nil {
		return err
	}
	defer f.Close()

	d, err := mp3.NewDecoder(f)
	if err != nil {
		return err
	}

	c, err := oto.NewContext(d.SampleRate(), 2, 2, 8192)
	if err != nil {
		return err
	}
	defer c.Close()

	p := c.NewPlayer()
	defer p.Close()

	fmt.Printf("Length: %d[bytes]\n", d.Length())

	if _, err := io.Copy(p, d); err != nil {
		return err
	}
	return nil
}

func main() {
	http.HandleFunc("/", viewHandler)
	http.HandleFunc("/api/attention", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
