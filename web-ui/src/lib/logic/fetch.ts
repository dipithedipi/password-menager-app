
// send the request with body and cookie
async function waitfetchData(url: string, method:string, body:any): Promise<{data:any, success:any}> {
    try {
        const response = await fetch(url, {
            method: method,
            credentials: 'include',
            headers: {
                'Content-Type': 'application/json',
                'Access-Control-Allow-Credentials': 'true',
            },
            body: method === "GET" ? null : JSON.stringify(body)
        });

        if (!response.ok) {
            throw new Error(`HTTP error status: ${response.status}`);
        }

        const data = await response.json();
        return {data, success: true};
    } catch (error: any) {
        console.error('Error fetching data:', error);
        return {data:null, success: false};
    }
}

export { waitfetchData };