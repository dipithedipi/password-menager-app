import { goto } from "$app/navigation";

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

        const data = await response.json();

        if (response.status === 401 && data.message.includes("Token")) {
            console.error("Token expired, redirect to login page");
            window.location.href = "/login";
            return {data: null, success: false};
        } else if (!response.ok) {
            throw new Error(data.message);
        }

        return {data, success: true};
    } catch (error: any) {
        console.error('Error fetching data:', error);
        return {data:error, success: false};
    }
}

async function getCategory() {
    let {data, success} = await waitfetchData('http://127.0.0.1:8000/category/gets', 'GET', {});
    if (!success) {
        console.error(data);
        goto("/login");
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

async function getEvents(startFormatted: string, endFormatted: string) {
    let {data, success} = await waitfetchData("http://127.0.0.1:8000/event/gets", "POST", {
			start: startFormatted, 
			end: endFormatted}
		);
		return {data, success};
}

export { waitfetchData, getCategory, resolveCategories, getEvents};