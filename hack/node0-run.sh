# get shell path
set -euo pipefail
SOURCE="$0"
while [ -h "$SOURCE"  ]; do
    DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
    SOURCE="$(readlink "$SOURCE")"
    [[ $SOURCE != /*  ]] && SOURCE="$DIR/$SOURCE"
done
DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
cd $DIR/../

# ensure host /chain_data points to hack/node0/chain_data (cannot change enclave.json)
ROOT_DIR="$(pwd)"
NODE_DIR="$ROOT_DIR/hack/node0"
NODE_CHAIN_DIR="$NODE_DIR/chain_data"
mkdir -p "$NODE_CHAIN_DIR"

ensure_chain_data_link() {
    local desired="$1"
    local target="/chain_data"

    if [ -L "$target" ]; then
        local cur
        cur="$(readlink "$target")"
        if [ "$cur" != "$desired" ]; then
            echo "ERROR: $target is a symlink to '$cur', expected '$desired'." 1>&2
            return 1
        fi
        return 0
    fi

    if [ -e "$target" ]; then
        # If it's an actual directory/file, don't destroy user data silently.
        # For enclave.json hostfs mount, an existing directory also works,
        # but it won't satisfy the requested mapping to hack/node0.
        echo "ERROR: $target exists and is not a symlink; please remove/move it to allow linking to '$desired'." 1>&2
        return 1
    fi

    if ln -s "$desired" "$target" 2>/dev/null; then
        return 0
    fi

    # Non-interactive environments can't prompt for sudo password. Fail with guidance.
    if command -v sudo >/dev/null 2>&1; then
        echo "ERROR: cannot create $target without privileges (and sudo needs a password)." 1>&2
        echo "Fix (one-time): sudo mkdir -p /chain_data && sudo rm -f /chain_data && sudo ln -s \"$desired\" /chain_data" 1>&2
        echo "Or: sudo mkdir -p /chain_data && sudo chown -R \"$(id -un)\" /chain_data" 1>&2
        return 1
    fi

    echo "ERROR: cannot create $target (no permission, and sudo not available)." 1>&2
    return 1
}

ensure_chain_data_link "$NODE_CHAIN_DIR" || exit 1

# build binary
ego-go build -o ./hack/build/dsecret ./main.go

# sign binary
cd ./hack/build/
ego sign dsecret

export SIDE_CHAIN_PORT=30110
export SIDE_CHAIN_RPC_PORT=30111
export GQL_PORT=30115
export CHAIN_ADDR=ws://192.168.110.205:30002/ws
# run from / so relative ./chain_data resolves to /chain_data
cd /
ego run "$ROOT_DIR/hack/build/dsecret"