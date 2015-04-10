from fabric.api import local, cd, put, env, shell_env, run
from fabric.context_managers import lcd
from fabric.contrib.project import rsync_project

env.use_ssh_config = True
env.hosts = ['root@server.artquest.ninja']

def deploy(what='all'):

    with shell_env(GOOS="linux", GOARCH="amd64"):
        local('godep go build')

    with lcd('frontend'):
        local('gulp build')

    with cd('/srv/artquest'):
        run("stop artquest || true")
        rsync_project("/srv/artquest/static", "frontend/dist/")
        put('rollingballs-server', 'artquest-server', mode=0755)
        run("start artquest")
