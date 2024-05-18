<script lang="ts">
    import { Modal } from 'flowbite';

    import { masterPassword } from '$lib/store/passwordStore';
    import { decryptAES } from '$lib/logic/cryptography';
	import MyModalNewPassword from './../../lib/components/MyModalNewPassword.svelte';
	import MyModalPasswordInfo from '$lib/components/MyModalPasswordInfo.svelte';
	import PasswordLine from '$lib/components/PasswordLine.svelte';

    import { getCategory } from '$lib/logic/fetch';
    import { waitfetchData } from '$lib/logic/fetch';
    import { onMount } from 'svelte';

    let passwordTitle = '';
    let passwordUsername = '';
    let passwordDescription = '';
    let passwordCategory = '';
    let passwordOtpProtected = false;

    let passwords:any = [];

    let searchBarValue = '';
    async function search(e: any, searchValue: string = searchBarValue, categories: string[] = categoriesSelected) {
        console.log(searchBarValue)
        if (searchValue === '') searchValue = '*';
        if (searchValue.trim() === '') return;
        let {data, success} = await waitfetchData('http://127.0.0.1:8000/password/search', 'POST', {
            domain: searchValue,
            category: categories
        });
        if (!success) {
            console.error(data);
            passwords = [];
            return;
        }
        passwords = data.passwords;
    }

    let categories:any[] = [];
    async function getCategoryHandle() {
        categories = await getCategory();
        categoriesSelected = categories.map((cat: any) => cat.name);
        console.log(categories);
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

    let categoriesSelected: string[] = [];
    function printSelectedCategories() {
        console.log(categoriesSelected);
    }

    function handleClickFilterCategory(e: any) {
        let category = e.target.value;

        if (categoriesSelected.includes(category)) {
            // one category is always selected
            if (categoriesSelected.length === 1) {
                alert("At least one category must be selected");
                // recheck the checkbox
                e.target.checked = true;
                return;
            }

            // remove category if it is already selected
            categoriesSelected = categoriesSelected.filter((cat) => cat !== category);
        } else {
            categoriesSelected.push(category);
        }
        // remove if a category is deselected 
        if (categoriesSelected.length === 0) {
            categoriesSelected = ["*"];
        }

        search(null, searchBarValue, categoriesSelected);
    }

    onMount(async () => {
        getCategoryHandle();
        search(null, "*", ["*"]);
    });

    function clickLineHandler(e: any, index: number, categories: any[]) {
        passwordTitle = passwords[index].website;
        passwordUsername = decryptAES(passwords[index].username, $masterPassword);
        passwordDescription = decryptAES(passwords[index].description, $masterPassword);
        passwordCategory = resolveCategories(passwords[index].category, categories);
        passwordOtpProtected = passwords[index].otpProtected;
        modal?.show();
    }

    function refreshPasswords() {
        console.log("Refreshing passwords");
        search(null, searchBarValue, categoriesSelected);
    }
    let modal: Modal | null = null;
</script>

<MyModalPasswordInfo title={passwordTitle} username={passwordUsername} description={passwordDescription === '' ? "No informations": passwordDescription} category={passwordCategory} otpProtected={passwordOtpProtected} modalId="passwordInfo" on:modalDetect={(e) => modal = e.detail}></MyModalPasswordInfo>
<MyModalNewPassword modalId="passwordAdd" on:updatePasswords={refreshPasswords}></MyModalNewPassword>

<section class="h-screen pb-2 mt-4 bg-gray-50 dark:bg-gray-700 pt-6 rounded-md">
    <div class="h-full mx-auto max-w-screen-xl px-4 lg:px-12">
        <div class="bg-white dark:bg-gray-800 relative shadow-md rounded-md sm:rounded-lg overflow-hidden">
            <div class="flex flex-col md:flex-row items-center justify-between space-y-3 md:space-y-0 md:space-x-4 p-4">
                <div class="w-full md:w-1/2">
                    <div class="flex items-center">
                        <label for="simple-search" class="sr-only">Search by website</label>
                        <div class="relative w-full">
                            <div class="absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none">
                                <svg aria-hidden="true" class="w-5 h-5 text-gray-500 dark:text-gray-400" fill="currentColor" viewbox="0 0 20 20" xmlns="http://www.w3.org/2000/svg">
                                    <path fill-rule="evenodd" d="M8 4a4 4 0 100 8 4 4 0 000-8zM2 8a6 6 0 1110.89 3.476l4.817 4.817a1 1 0 01-1.414 1.414l-4.816-4.816A6 6 0 012 8z" clip-rule="evenodd" />
                                </svg>
                            </div>
                            <input type="text" bind:value={searchBarValue} on:keyup={search} id="simple-search" class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-primary-500 focus:border-primary-500 block w-full pl-10 p-2 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-primary-500 dark:focus:border-primary-500" placeholder="Search by website">
                        </div>
                    </div>
                </div>
                <div class="w-full md:w-auto flex flex-col md:flex-row space-y-2 md:space-y-0 items-stretch md:items-center justify-end md:space-x-3 flex-shrink-0">
                    <button type="button" data-modal-target="passwordAdd" data-modal-toggle="passwordAdd" class="flex items-center justify-center text-white bg-primary-700 hover:bg-primary-800 focus:ring-4 focus:ring-primary-300 font-medium rounded-lg text-sm px-4 py-2 dark:bg-primary-600 dark:hover:bg-primary-700 focus:outline-none dark:focus:ring-primary-800">
                        <svg class="h-3.5 w-3.5 mr-2" fill="currentColor" viewbox="0 0 20 20" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
                            <path clip-rule="evenodd" fill-rule="evenodd" d="M10 3a1 1 0 011 1v5h5a1 1 0 110 2h-5v5a1 1 0 11-2 0v-5H4a1 1 0 110-2h5V4a1 1 0 011-1z" />
                        </svg>
                        Add password
                    </button>
                    <div class="flex items-center space-x-3 w-full md:w-auto">
                        <button id="filterDropdownButton" on:click={printSelectedCategories} data-dropdown-toggle="filterDropdown" class="w-full md:w-auto flex items-center justify-center py-2 px-4 text-sm font-medium text-gray-900 focus:outline-none bg-white rounded-lg border border-gray-200 hover:bg-gray-100 hover:text-primary-700 focus:z-10 focus:ring-4 focus:ring-gray-200 dark:focus:ring-gray-700 dark:bg-gray-800 dark:text-gray-400 dark:border-gray-600 dark:hover:text-white dark:hover:bg-gray-700" type="button">
                            <svg xmlns="http://www.w3.org/2000/svg" aria-hidden="true" class="h-4 w-4 mr-2 text-gray-400" viewbox="0 0 20 20" fill="currentColor">
                                <path fill-rule="evenodd" d="M3 3a1 1 0 011-1h12a1 1 0 011 1v3a1 1 0 01-.293.707L12 11.414V15a1 1 0 01-.293.707l-2 2A1 1 0 018 17v-5.586L3.293 6.707A1 1 0 013 6V3z" clip-rule="evenodd" />
                            </svg>
                            Filter
                            <svg class="-mr-1 ml-1.5 w-5 h-5" fill="currentColor" viewbox="0 0 20 20" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
                                <path clip-rule="evenodd" fill-rule="evenodd" d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z" />
                            </svg>
                        </button>
                        <div id="filterDropdown" class="z-10 border-2 border-gray-500 hidden w-48 p-3 bg-white rounded-lg shadow dark:bg-gray-700">
                            <h6 class="mb-3 text-sm font-medium text-gray-900 dark:text-white">Choose category</h6>
                            <ul class="space-y-2 text-sm" aria-labelledby="filterDropdownButton">
                                {#each categories as category}
                                    <li class="flex items-center">
                                        <input bind:value={category.name} on:click={handleClickFilterCategory} checked type="checkbox" class="w-4 h-4 bg-gray-100 border-gray-300 rounded text-primary-600 focus:ring-primary-500 dark:focus:ring-primary-600 dark:ring-offset-gray-700 focus:ring-2 dark:bg-gray-600 dark:border-gray-500">
                                        <label for="category" class="ml-2 text-sm font-medium text-gray-900 dark:text-gray-100">{category.name}</label>
                                    </li>
                                {/each}  
                            </ul>
                        </div>
                    </div>
                </div>
            </div>
            <div class="overflow-x-auto">
                <table class="w-full text-sm text-left text-gray-500 dark:text-gray-400">
                    <thead class="text-xs text-gray-700 uppercase bg-gray-50 dark:bg-gray-700 dark:text-gray-400">
                        <tr>
                            <th scope="col" class="px-4 py-3">Logo</th>
                            <th scope="col" class="px-4 py-3">Email</th>
                            <th scope="col" class="px-4 py-3 hidden sm:flex">Website</th>
                            <th scope="col" class="px-4 py-3">Last use</th>
                            <th scope="col" class="px-4 py-3">
                                <span class="sr-only">Actions</span>
                            </th>
                        </tr>
                    </thead>
                    <tbody>
                        <!-- repeat the component n times based on the array length -->
                        {#each passwords as password, index}
                            <PasswordLine username={decryptAES(password.username, $masterPassword)} website={password.website} lastuse={password.lastUsed} on:click={ (event) => clickLineHandler(event, index, categories) }></PasswordLine>
                        {/each}
                        {#if passwords.length === 0}
                            <tr>
                                <td class="px-4 py-3 text-center" colspan="5">No password found</td>
                            </tr>
                        {/if}
                    </tbody>
                </table>
            </div>
            <nav class="flex flex-col md:flex-row justify-between items-start md:items-center space-y-3 md:space-y-0 p-4" aria-label="Table navigation">
                <span class="text-sm font-normal text-gray-500 dark:text-gray-400">
                    Showing
                    <span class="font-semibold text-gray-900 dark:text-white">1-10</span>
                    of
                    <span class="font-semibold text-gray-900 dark:text-white">1000</span>
                </span>
                <ul class="inline-flex items-stretch -space-x-px">
                    <li>
                        <a href="#" class="flex items-center justify-center h-full py-1.5 px-3 ml-0 text-gray-500 bg-white rounded-l-lg border border-gray-300 hover:bg-gray-100 hover:text-gray-700 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-white">
                            <span class="sr-only">Previous</span>
                            <svg class="w-5 h-5" aria-hidden="true" fill="currentColor" viewbox="0 0 20 20" xmlns="http://www.w3.org/2000/svg">
                                <path fill-rule="evenodd" d="M12.707 5.293a1 1 0 010 1.414L9.414 10l3.293 3.293a1 1 0 01-1.414 1.414l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 0z" clip-rule="evenodd" />
                            </svg>
                        </a>
                    </li>
                    <li>
                        <a href="#" class="flex items-center justify-center text-sm py-2 px-3 leading-tight text-gray-500 bg-white border border-gray-300 hover:bg-gray-100 hover:text-gray-700 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-white">1</a>
                    </li>
                    <li>
                        <a href="#" class="flex items-center justify-center text-sm py-2 px-3 leading-tight text-gray-500 bg-white border border-gray-300 hover:bg-gray-100 hover:text-gray-700 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-white">2</a>
                    </li>
                    <li>
                        <a href="#" aria-current="page" class="flex items-center justify-center text-sm z-10 py-2 px-3 leading-tight text-primary-600 bg-primary-50 border border-primary-300 hover:bg-primary-100 hover:text-primary-700 dark:border-gray-700 dark:bg-gray-700 dark:text-white">3</a>
                    </li>
                    <li>
                        <a href="#" class="flex items-center justify-center text-sm py-2 px-3 leading-tight text-gray-500 bg-white border border-gray-300 hover:bg-gray-100 hover:text-gray-700 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-white">...</a>
                    </li>
                    <li>
                        <a href="#" class="flex items-center justify-center text-sm py-2 px-3 leading-tight text-gray-500 bg-white border border-gray-300 hover:bg-gray-100 hover:text-gray-700 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-white">100</a>
                    </li>
                    <li>
                        <a href="#" class="flex items-center justify-center h-full py-1.5 px-3 leading-tight text-gray-500 bg-white rounded-r-lg border border-gray-300 hover:bg-gray-100 hover:text-gray-700 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-white">
                            <span class="sr-only">Next</span>
                            <svg class="w-5 h-5" aria-hidden="true" fill="currentColor" viewbox="0 0 20 20" xmlns="http://www.w3.org/2000/svg">
                                <path fill-rule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clip-rule="evenodd" />
                            </svg>
                        </a>
                    </li>
                </ul>
            </nav>
        </div>
    </div>
</section>