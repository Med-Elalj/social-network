export async function SendData(url, Data) {
    try {
        const response = await fetch(url, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(Data),
        })
        
        const responseBody = await response.json()
        return { status: response.status, body: responseBody }
    } catch (error) {
        return
    }
}