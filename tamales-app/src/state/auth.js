import EventSource from 'react-native-event-source'

export const waitForAuthToken = (email) => new Promise((resolve) => {
    console.log('waitForAuthToken')
    const body = JSON.stringify({email});
    const options = {method: 'POST', headers: {Accept: 'text/event-stream'}, body}
    const eventSource = new EventSource('http://tamales.eighty4.io/login/initiate', options)
    eventSource.addEventListener('success', event => {
        console.log('data')
        console.log(typeof event.data, event.data)
        resolve(event.data)
    })
})
