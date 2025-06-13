export async function SendData(url, Data) {   
    try {
        const response = await fetch(url, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(Data),
        })
        
        return response
    } catch (error) {
        return        
    }
}