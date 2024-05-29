<script lang="ts">
	import { masterPassword } from '$lib/store/passwordStore';
	import { decryptAES, hashSHA256 } from '$lib/logic/cryptography';
	import { waitfetchData } from '$lib/logic/fetch';
    import { onMount } from 'svelte';
    
    let compromised = 0;
    let same = 0;
    let totalPassword = 0;

    let loading = true;

    let compromisedPasswords:any = [];
    let samePasswords:any = []

    async function getPossibleProblems() {
        console.log('Health check');
        let {data, success} = await waitfetchData('http://127.0.0.1:8000/password/health', 'GET', {});
        
        if (!success) {
            console.log('Failed to fetch data');
            return;
        }

        compromised = data.compromisedPasswordCount;
        same = data.samePasswordCount;
        totalPassword = data.totalPassword;
        compromisedPasswords = data.compromisedPassword;
        samePasswords = data.samePassword;
    }

    function checkFalseResultCompromised() {
        if (compromisedPasswords===null || compromisedPasswords.length == 0) {
            return
        }

        compromisedPasswords.forEach((compromisedPassword: any) => {
            //console.log("password:", compromisedPassword.password);
            if (hashSHA256(decryptAES(compromisedPassword.password, $masterPassword)) !== compromisedPassword.LeakedHash) {
                //console.log("triggered");
                compromised -= 1;
                compromisedPasswords = compromisedPasswords.filter((item: any) => item !== compromisedPassword)
            }
            // console.log("password:", decryptAES(compromisedPassword.password, $masterPassword), 
            // "hash:", hashSHA256(decryptAES(compromisedPassword.password, $masterPassword)), 
            // "partial hash:", compromisedPassword.Partialhash,
            // "leaked hash:", compromisedPassword.LeakedHash);
        });
    }

    function checkFalseResultSame(): void {
        if (samePasswords===null || samePasswords.length === 0) {
            return;
        }

        // flat the nested array
        let samePasswordsFlatted = samePasswords.flat()
        console.log(samePasswordsFlatted)

        samePasswordsFlatted.forEach((samePassword: any) => {
            samePassword.fullHash = hashSHA256(decryptAES(samePassword.password, $masterPassword))
        })

        // remove the hashes that are unique
        let notUniqueHash = samePasswordsFlatted.filter((item: any) => samePasswordsFlatted.filter((item2: any) => item2.fullHash === item.fullHash).length > 1)

        // split the array in sub array where the hash is the same
        let samePasswordsGrouped = notUniqueHash.reduce((acc: any, item: any) => {
            let key = item.fullHash;
            if (!acc[key]) {
                acc[key] = [];
            }
            acc[key].push(item);
            return acc;
        }, {});

        samePasswords = Object.values(samePasswordsGrouped)
    }

    async function healthCheck() {
        compromisedPasswords= [];
        samePasswords = []
        compromised = 0;
        same = 0;
        totalPassword = 0;
        loading = true;

        await getPossibleProblems();
        checkFalseResultCompromised();
        checkFalseResultSame();
        loading = false;
    }

    onMount(async () => {
        await healthCheck();
    });
</script>

<div class="h-screen bg-gray-700">
    <section class="pt-6 dark:bg-gray-700 flex flex-col justify-center items-center">
        <div class="max-w-screen-xl px-4 py-4 mx-auto text-center lg:px-6">
            <dl class="grid max-w-screen-md gap-7 mx-auto text-gray-900 sm:grid-cols-3 dark:text-white">
                <div class="flex flex-col dark:bg-gray-800 items-center justify-center p-2 border-2 border-dashed border-slate-500 rounded-md">
                    <dt class="mb-2 text-3xl md:text-4xl font-extrabold">{totalPassword}</dt>
                    <dd class="font-light text-gray-500 dark:text-gray-400">Total password</dd>
                </div>
                <div class="flex flex-col items-center dark:bg-gray-800 justify-center p-2 text-yellow-300 border-2 border-dashed border-slate-500 rounded-md">
                    <dt class="mb-2 text-3xl md:text-4xl font-extrabold">{same}</dt>
                    <dd class="font-light text-gray-500 dark:text-gray-400">Reused password</dd>
                </div>
                <div class="flex flex-col items-center dark:bg-gray-800 justify-center p-2 text-red-600 border-2 border-dashed border-slate-500 rounded-md">
                    <dt class="mb-2 text-3xl md:text-4xl font-extrabold">{compromised}</dt>
                    <dd class="font-light text-gray-500 dark:text-gray-400">Compromised</dd>
                </div>

                <div class="pl-2 flex items-center justify-center col-span-3">
                    <button on:click={healthCheck} type="button" class="px-3 py-2 mb-3 w-full mr-2 sm:mr:1 text-sm font-medium text-center text-white bg-blue-800 rounded-lg hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-400 dark:bg-blue-700 dark:hover:bg-blue-700 dark:focus:ring-blue-800">
                        Check all passwords now
                    </button>
                </div>
            </dl>
        </div>
        {#if loading}
            <div class="flex items-center justify-center w-48 h-48 rounded-lg  dark:bg-gray-700">
                <div role="status">
                    <svg aria-hidden="true" class="w-8 h-8 text-gray-200 animate-spin dark:text-gray-600 fill-blue-600" viewBox="0 0 100 101" fill="none" xmlns="http://www.w3.org/2000/svg"><path d="M100 50.5908C100 78.2051 77.6142 100.591 50 100.591C22.3858 100.591 0 78.2051 0 50.5908C0 22.9766 22.3858 0.59082 50 0.59082C77.6142 0.59082 100 22.9766 100 50.5908ZM9.08144 50.5908C9.08144 73.1895 27.4013 91.5094 50 91.5094C72.5987 91.5094 90.9186 73.1895 90.9186 50.5908C90.9186 27.9921 72.5987 9.67226 50 9.67226C27.4013 9.67226 9.08144 27.9921 9.08144 50.5908Z" fill="currentColor"/><path d="M93.9676 39.0409C96.393 38.4038 97.8624 35.9116 97.0079 33.5539C95.2932 28.8227 92.871 24.3692 89.8167 20.348C85.8452 15.1192 80.8826 10.7238 75.2124 7.41289C69.5422 4.10194 63.2754 1.94025 56.7698 1.05124C51.7666 0.367541 46.6976 0.446843 41.7345 1.27873C39.2613 1.69328 37.813 4.19778 38.4501 6.62326C39.0873 9.04874 41.5694 10.4717 44.0505 10.1071C47.8511 9.54855 51.7191 9.52689 55.5402 10.0491C60.8642 10.7766 65.9928 12.5457 70.6331 15.2552C75.2735 17.9648 79.3347 21.5619 82.5849 25.841C84.9175 28.9121 86.7997 32.2913 88.1811 35.8758C89.083 38.2158 91.5421 39.6781 93.9676 39.0409Z" fill="currentFill"/></svg>
                    <span class="sr-only">Loading...</span>
                </div>
            </div>
        {:else}
            <div>
                {#if compromised === 0}
                    <h2 class="mb-2 mt-2 text-lg font-semibold text-gray-900 dark:text-white">No compromised password found</h2>
                {:else}
                    <h2 class="mb-2 text-lg font-semibold text-gray-900 dark:text-white">Compromised:</h2>
                    {#each compromisedPasswords as compromisedPassword}
                        <ul class="max-w-md space-y-1 pb-2 text-gray-500 list-disc list-inside dark:text-gray-400">
                            <li class="flex dark:bg-gray-800 items-center border-2 rounded-lg border-red-500 p-4">
                                <svg class="scale-125 w-3.5 h-3.5 me-2 text-red-500 dark:text-red-400 flex-shrink-0" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="currentColor" viewBox="0 0 24 24">
                                    <path d="M12.8659 3.00017L22.3922 19.5002C22.6684 19.9785 22.5045 20.5901 22.0262 20.8662C21.8742 20.954 21.7017 21.0002 21.5262 21.0002H2.47363C1.92135 21.0002 1.47363 20.5525 1.47363 20.0002C1.47363 19.8246 1.51984 19.6522 1.60761 19.5002L11.1339 3.00017C11.41 2.52187 12.0216 2.358 12.4999 2.63414C12.6519 2.72191 12.7782 2.84815 12.8659 3.00017ZM4.20568 19.0002H19.7941L11.9999 5.50017L4.20568 19.0002ZM10.9999 16.0002H12.9999V18.0002H10.9999V16.0002ZM10.9999 9.00017H12.9999V14.0002H10.9999V9.00017Z"></path>
                                </svg>
                                {decryptAES(compromisedPassword.username, $masterPassword)} for {compromisedPassword.website}
                            </li>
                        </ul>
                    {/each}
                {/if}
            </div>
            <div class="pt-6">
                {#if same === 0}
                    <h2 class="mb-2 mt-2 text-lg font-semibold text-gray-900 dark:text-white">No reused password found</h2>
                {:else}
                    {#each samePasswords as samePassword}
                    <h2 class="mb-2 mt-2 text-lg font-semibold text-gray-900 dark:text-white">Same passwords:</h2>
                        {#each samePassword as singlePassword}
                            <ul class="max-w-md space-y-1 pb-2 text-gray-500 list-disc list-inside dark:text-gray-400">
                                <li class="flex dark:bg-gray-800 items-center border-2 rounded-lg border-yellow-300 p-4">
                                    <svg class="scale-125 w-3.5 h-3.5 me-2 text-yellow-500 dark:text-yellow-400 flex-shrink-0" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="currentColor" viewBox="0 0 24 24">
                                        <path d="M12.8659 3.00017L22.3922 19.5002C22.6684 19.9785 22.5045 20.5901 22.0262 20.8662C21.8742 20.954 21.7017 21.0002 21.5262 21.0002H2.47363C1.92135 21.0002 1.47363 20.5525 1.47363 20.0002C1.47363 19.8246 1.51984 19.6522 1.60761 19.5002L11.1339 3.00017C11.41 2.52187 12.0216 2.358 12.4999 2.63414C12.6519 2.72191 12.7782 2.84815 12.8659 3.00017ZM4.20568 19.0002H19.7941L11.9999 5.50017L4.20568 19.0002ZM10.9999 16.0002H12.9999V18.0002H10.9999V16.0002ZM10.9999 9.00017H12.9999V14.0002H10.9999V9.00017Z"></path>
                                    </svg>
                                    {decryptAES(singlePassword.username, $masterPassword)} for {singlePassword.website}
                                </li>
                            </ul>
                        {/each}
                    {/each}
                {/if}
            </div>
        {/if}
    </section>
</div>