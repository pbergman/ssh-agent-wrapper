build:
	go build -o ssh-agent-wrapper -x -ldflags "-s -w" && strip ssh-agent-wrapper
install: ssh-agent-wrapper
	cp ssh-agent-wrapper /usr/local/bin
clean: ssh-agent-wrapper
	rm ssh-agent-wrapper
uninstall: /usr/local/bin/ssh-agent-wrapper
	rm 	/usr/local/bin/ssh-agent-wrapper
