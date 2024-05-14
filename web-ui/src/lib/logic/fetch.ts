async function waitfetchData(url: string, method:string, body:any): Promise<any> {
    try {
        const response = await fetch(url, {
            method: method,
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(body),
        });

        if (!response.ok) {
            throw new Error(`HTTP error status: ${response.status}`);
        }

        const data = await response.json();
        return data;
    } catch (error: any) {
        console.error('Error fetching data:', error);

        // check if the error is because the token is expired and redirect to login
        if (error.message.contains('Token')) {
            alert(`${error.message}, Redirecting to login page.`);
            window.location.href = '/login';
        }
        return false;
    }
}

export { waitfetchData };