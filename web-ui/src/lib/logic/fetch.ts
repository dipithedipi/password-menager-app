
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

async function getCategory() {
    let {data, success} = await waitfetchData('http://127.0.0.1:8000/category/gets', 'GET', {});
    if (!success) {
        console.error(data);
        return;
    }
    return data.categories;
}

function resolveCategories(category: string, categories: any[]) {
    console.log("Categories data", categories)
    if (!categories ||!Array.isArray(categories)) {
        console.log("Categories data not available");
        return;
    }
    let categoryObj = categories.find((cat: any) => cat.id === category);
    return categoryObj?.name;
}

export { waitfetchData, getCategory, resolveCategories};