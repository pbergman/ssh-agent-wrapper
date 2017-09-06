ssh-agent-wrapper is a simple wrapper around `ssh-agent` and `ssh-add`
for starting the ssh-agent and adding keys (on start). It will check
if the "ssh agent env" file exist and if not then it will start the
`ssh-agent` process and save the output to the ssh agent env file.
That can later be used to determine pid etc. of the ssh-agent process.

The file can be given with the `-f` or `--file` option and haves a
default of ~/.ssh/agent_env (see also `-h`)

To run this script with the login of a shell you can add the following
line to .bashrc, .profile or any other startup script. (eval is used for
setting the system envs passed by the ssh-agent or read from envs file)

```
eval `ssh-agent-wrapper`
```	

To install ssh-agent-wrapper you can do:

```
make
sudo make install
```