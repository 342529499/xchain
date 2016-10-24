target1=vagrant@192.168.1.230
target2=vagrant@192.168.1.231


function releaseOnServer() {
    ip=$1

    cd $GOPATH/src/github.com/1851616111/xchain
    git pull
    go build
    ./xchain start --entryPointAddress=192.168.1.186:10690 --netAddress=$ip

}


#function release() {
#
#
#}

releaseOnServer 192.168.1.230

function main() {


}