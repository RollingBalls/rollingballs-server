from fabric.api import local, cd, put, env, shell_env, run

env.use_ssh_config = True
env.hosts = ['root@server.artquest.ninja']

def deploy(what='all'):

    with shell_env(GOOS="linux", GOARCH="amd64"):
        local('godep go build')

    with cd('/srv/artquest'):
        run("stop artquest")
        put('rollingballs-server', 'artquest-server', mode=0755)
        run("start artquest")
