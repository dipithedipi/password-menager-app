<script lang="ts">
    import { waitfetchData } from "$lib/logic/fetch";
    import { Modal } from "flowbite";
    import { createEventDispatcher } from "svelte";

    let code1: string;
    let code2: string;
    let code3: string;
    let code4: string;
    let code5: string;
    let code6: string;
    let otpCodeInput: boolean = false;
    let errorOtp: boolean = false;
    let errorOtpText: string = "";

    export let modalId: string;
    export let username: string;
    export let title: string;
    export let description: string;
    export let category: string;
    export let otpProtected: boolean;

    let modalElem: HTMLDivElement | null = null;
    let modal: Modal | null = null;
    $: if (modalElem!== null) {
        modal = new Modal(modalElem);
    }

    const dispatch = createEventDispatcher<{modalDetect: Modal}>();

    $: if (modal!== null) {
        dispatch("modalDetect", modal);
    }

    function getOtpCode() {
        return `${code1}${code2}${code3}${code4}${code5}${code6}`;
    }

    function deletePassword() {
        console.log("Delete password");
    }

    function copyPassword() {
        console.log("Copy password");
    }

    function editPassword() {
        console.log("Edit password");
    }
</script>

<div bind:this={modalElem} id={modalId} data-modal-target="{modalId}" tabindex="-1" aria-hidden="true" class="fixed hidden top-1/2 sm:top-0 -translate-y-1/4 sm:-translate-y-0 overflow-y-auto overflow-x-hidden right-0 left-0 z-50 justify-center items-center w-full md:inset-0 h-modal md:h-full">
    <div class="relative p-4 w-full max-w-xl h-full md:h-auto">
        <!-- Modal content -->
        <div class="relative p-4 bg-white rounded-lg shadow dark:bg-gray-800 sm:p-5">
                <!-- Modal header -->
                <div class="flex justify-between mb-4 rounded-t sm:mb-5">
                    <div class="text-lg text-gray-900 md:text-xl dark:text-white">
                        <h3 class="font-semibold">
                            {!otpCodeInput ? title: `Delete password: ${title}`}
                        </h3>                       
                    </div>
                    <div>
                        <button type="button" class="text-gray-400 bg-transparent hover:bg-gray-200 hover:text-gray-900 rounded-lg text-sm p-1.5 inline-flex dark:hover:bg-gray-600 dark:hover:text-white" on:click={() => modal?.hide()}>
                            <svg aria-hidden="true" class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg"><path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd"></path></svg>
                            <span class="sr-only">Close modal</span>
                        </button>
                    </div>
                </div>
                {#if !otpCodeInput}
                    <dl>
                        <dt class="mb-2 font-semibold leading-none text-gray-900 dark:text-white">{username.includes('@') ? "Email": "Username"}</dt>
                        <dd class="mb-4 font-light text-gray-500 sm:mb-5 dark:text-gray-400">{username}</dd>
                        <dt class="mb-2 font-semibold leading-none text-gray-900 dark:text-white">Details</dt>
                        <dd class="mb-4 font-light text-gray-500 sm:mb-5 dark:text-gray-400">{description}</dd>
                        <dt class="mb-2 font-semibold leading-none text-gray-900 dark:text-white">Category</dt>
                        <dd class="mb-4 font-light text-gray-500 sm:mb-5 dark:text-gray-400">{category}</dd>
                    </dl>
                {:else}
                    <dl>
                        <div class="max-w-sm mx-auto" >
                            <div class="flex mb-2 space-x-2 rtl:space-x-reverse mx-auto">
                                <div>
                                    <label for="code-1" class="sr-only">First code</label>
                                    <input bind:value={code1} type="text" maxlength="1" data-focus-input-init data-focus-input-next="code-2" id="code-1" class="block w-10 h-10 md:w-14 md:h-14  py-3 text-sm font-extrabold text-center text-gray-900 bg-white border border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-primary-500 dark:focus:border-primary-500"  />
                                </div>
                                <div>
                                    <label for="code-2" class="sr-only">Second code</label>
                                    <input bind:value={code2} type="text" maxlength="1" data-focus-input-init data-focus-input-prev="code-1" data-focus-input-next="code-3" id="code-2" class="block w-10 h-10 md:w-14 md:h-14  py-3 text-sm font-extrabold text-center text-gray-900 bg-white border border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-primary-500 dark:focus:border-primary-500"  />
                                </div>
                                <div>
                                    <label for="code-3" class="sr-only">Third code</label>
                                    <input bind:value={code3} type="text" maxlength="1" data-focus-input-init data-focus-input-prev="code-2" data-focus-input-next="code-4" id="code-3" class="block w-10 h-10 md:w-14 md:h-14  py-3 text-sm font-extrabold text-center text-gray-900 bg-white border border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-primary-500 dark:focus:border-primary-500"  />
                                </div>
                                <div>
                                    <label for="code-4" class="sr-only">Fourth code</label>
                                    <input bind:value={code4} type="text" maxlength="1" data-focus-input-init data-focus-input-prev="code-3" data-focus-input-next="code-5" id="code-4" class="block w-10 h-10 md:w-14 md:h-14  py-3 text-sm font-extrabold text-center text-gray-900 bg-white border border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-primary-500 dark:focus:border-primary-500"  />
                                </div>
                                <div>
                                    <label for="code-5" class="sr-only">Fifth code</label>
                                    <input bind:value={code5} type="text" maxlength="1" data-focus-input-init data-focus-input-prev="code-4" data-focus-input-next="code-6" id="code-5" class="block w-10 h-10 md:w-14 md:h-14  py-3 text-sm font-extrabold text-center text-gray-900 bg-white border border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-primary-500 dark:focus:border-primary-500"  />
                                </div>
                                <div>
                                    <label for="code-6" class="sr-only">Sixth code</label>
                                    <input bind:value={code6} type="text" maxlength="1" data-focus-input-init data-focus-input-prev="code-5" id="code-6" class="block w-10 h-10 md:w-14 md:h-14  py-3 text-sm font-extrabold text-center text-gray-900 bg-white border border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-primary-500 dark:focus:border-primary-500"  />
                                </div>
                            </div>
                            <p id="helper-text-explanation" class="mt-2 mb-4 text-sm text-gray-500 dark:text-gray-400">Please introduce the 6 digit code from Autenticator App.</p>
                            {#if errorOtp}
                                <p class="text-sm mb-2 text-red-600 dark:text-red-500">{errorOtpText}</p>
                            {/if}
                        </div>
                    </dl>
                {/if}
                <div class="flex justify-between items-center">
                    <div class="flex items-center space-x-5 sm:space-x-4">
                        <button on:click={editPassword} type="button" class="text-white inline-flex items-center bg-primary-700 hover:bg-primary-800 focus:ring-4 focus:outline-none focus:ring-primary-300 font-medium rounded-lg text-sm px-5 py-2.5 text-center dark:bg-primary-600 dark:hover:bg-primary-700 dark:focus:ring-primary-800">
                            <svg aria-hidden="true" class="mr-1 -ml-1 w-5 h-5" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg"><path d="M17.414 2.586a2 2 0 00-2.828 0L7 10.172V13h2.828l7.586-7.586a2 2 0 000-2.828z"></path><path fill-rule="evenodd" d="M2 6a2 2 0 012-2h4a1 1 0 010 2H4v10h10v-4a1 1 0 112 0v4a2 2 0 01-2 2H4a2 2 0 01-2-2V6z" clip-rule="evenodd"></path></svg>
                            Edit
                        </button>               
                        <button on:click={copyPassword} type="button" class="py-2.5 px-5 inline-flex text-sm font-medium text-gray-900 focus:outline-none bg-white rounded-lg border border-gray-200 hover:bg-gray-100 hover:text-primary-700 focus:z-10 focus:ring-4 focus:ring-gray-200 dark:focus:ring-gray-700 dark:bg-gray-800 dark:text-gray-400 dark:border-gray-600 dark:hover:text-white dark:hover:bg-gray-700">
                            <svg xmlns="http://www.w3.org/2000/svg" class="mr-1 -ml-1 w-5 h-5" viewBox="0 0 24 24" fill="currentColor"><path d="M7 4V2H17V4H20.0066C20.5552 4 21 4.44495 21 4.9934V21.0066C21 21.5552 20.5551 22 20.0066 22H3.9934C3.44476 22 3 21.5551 3 21.0066V4.9934C3 4.44476 3.44495 4 3.9934 4H7ZM7 6H5V20H19V6H17V8H7V6ZM9 4V6H15V4H9Z"></path></svg>
                            Copy
                            <!-- space -->
                            <div class="hidden ml-1 sm:block">
                                the password
                            </div>
                        </button>
                    </div>
                    <button type="button" on:click={deletePassword} class="inline-flex items-center text-white bg-red-600 hover:bg-red-700 focus:ring-4 focus:outline-none focus:ring-red-300 font-medium rounded-lg text-sm px-5 py-2.5 text-center dark:bg-red-500 dark:hover:bg-red-600 dark:focus:ring-red-900">
                        <svg aria-hidden="true" class="w-5 h-5 sm:mr-1.5 sm:-ml-1" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg"><path fill-rule="evenodd" d="M9 2a1 1 0 00-.894.553L7.382 4H4a1 1 0 000 2v10a2 2 0 002 2h8a2 2 0 002-2V6a1 1 0 100-2h-3.382l-.724-1.447A1 1 0 0011 2H9zM7 8a1 1 0 012 0v6a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v6a1 1 0 102 0V8a1 1 0 00-1-1z" clip-rule="evenodd"></path></svg>
                        <div>
                            Delete
                        </div>
                    </button> 
                </div>
        </div>
    </div>
</div>