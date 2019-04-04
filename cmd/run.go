package cmd

import (
	"filemux/config"
	"filemux/processor"
	"filemux/processor/broadcaster"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	logger "github.com/apsdehal/go-logger"
	figure "github.com/common-nighthawk/go-figure"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
)

var (
	wg sync.WaitGroup

	runCmd = &cobra.Command{
		Use:   "run",
		Short: "Starts the watchers then blocks",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(figure.NewFigure(AppName, "", true).String())

			log.SetLogLevel(logger.LogLevel(logLevel))

			conf, err := initConfig(confFilePath)
			if err != nil {
				log.Fatal(err.Error())
			}

			br := broadcaster.New()
			go br.Run()

			processors := conf.Processors()

			for _, p := range processors {
				if err := p.Start(br.Pub); err != nil {
					log.Fatal(err.Error())
				}
			}

			sigs := make(chan os.Signal, 1)

			signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

			go func() {
				s := <-sigs

				log.Noticef("Received %s: gracefully stopping everything....", s)

				for _, p := range processors {
					if err := p.Stop(); err != nil {
						log.Criticalf("An error occurred while stopping %s: %s", p, err)
					}
				}

				wg.Wait()

				os.Exit(0)
			}()

			if serverListener != "" {
				go func() {
					log.Fatal(http.ListenAndServe(serverListener, router(conf, br)).Error())
				}()
			}

			brRecv, _ := br.Recv()

			for {
				select {
				case e := <-brRecv:
					log.Log(e.Level, fmt.Sprintf("%s %s", e.Processor, e.Message))

					switch e.Type {
					case processor.Bgn:
						wg.Add(1)
					case processor.Fin:
						wg.Done()
					}
				}
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringVarP(&confFilePath, "conf-file", "c", "", "configuration file path")
	runCmd.MarkFlagRequired("conf-file")

	runCmd.Flags().StringVarP(&serverListener, "api", "a", "", "listen for incomming HTTP/WebSocket requests on x:x (e.g. :8080 or 127.0.0.1:8080)")

	runCmd.Flags().Var(&logLevel, "log-level", "")

	log, _ = logger.New(AppName, 0, os.Stdout)

	log.SetFormat("[%{time}] [%{level}] %{message}")
}

func initConfig(confFilePath string) (*config.Conf, error) {
	confFileContent, err := ioutil.ReadFile(confFilePath)
	if err != nil {
		return nil, fmt.Errorf("Cannot read the config file: %s", err)
	}

	conf, err := config.New(confFileContent)
	if err != nil {
		return nil, err
	}

	return conf, nil
}

func router(conf *config.Conf, broadcaster *broadcaster.Broadcaster) *mux.Router {
	r := mux.NewRouter().StrictSlash(true)

	wsUpgrader := websocket.Upgrader{}

	r.PathPrefix("/ws/events").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := wsUpgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Error(err.Error())

			return
		}

		eventsRecv, eventsRecvUUID := broadcaster.Recv()
		tk := time.NewTicker(time.Second)

		defer func() {
			c.Close()

			broadcaster.RemoveRecv(eventsRecvUUID)
			tk.Stop()
		}()

		for _, p := range conf.Processors() {
			if err := c.WriteJSON(processor.Event{Processor: p, Type: processor.Log, Message: "", Level: logger.NoticeLevel}); err != nil {
				log.Info(err.Error())

				return
			}
		}

		done := make(chan int)
		defer func() {
			done <- 1
		}()

		go func() {
			for {
				select {
				case e, ok := <-eventsRecv:
					if ok {
						if err := c.WriteJSON(e); err != nil {
							log.Info(err.Error())

							return
						}
					}
				case <-tk.C:
					if err := c.WriteMessage(websocket.PingMessage, []byte("ping")); err != nil {
						log.Info(err.Error())

						return
					}
				case <-done:
					return
				}
			}
		}()

		for {
			_, _, err := c.ReadMessage()
			if err != nil {
				return
			}
		}
	}).Name("events")

	return r
}
