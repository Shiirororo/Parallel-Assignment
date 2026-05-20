import http from 'k6/http'

export let GetInfo_options = {
    stages: [
        { duration: '120s', target: 100 },
        { duration: '90s', target: 300 },
        { duration: '90s', target: 800 },
        { duration: '120s', target: 1000 },

    ]
}

export let Register_options = {
    stages: [
        { duration: '120s', target: 100 },

    ]
}

export let Unregister_options = {
    stages: [
        { duration: '120s', target: 100 },

    ]
}
export default function () {
    http.get()
}