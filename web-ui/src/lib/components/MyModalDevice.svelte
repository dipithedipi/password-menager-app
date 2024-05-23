<script lang="ts">
	import { waitfetchData } from "$lib/logic/fetch";
    import { otpCodeInputValue } from "$lib/store/otpStore";
    import { betterTime } from "$lib/logic/utils";
    import { createEventDispatcher } from "svelte";
    import OtpInput from "./OtpInput.svelte";

    const dispatch = createEventDispatcher();

    let otpInputStatus = false;
    let errorOtp = false;

    export let modalId: string;
    export let modal: any;

    export let deviceId: string;
    export let activeSession: boolean;
    export let userAgent: string;
    export let userAgentFormatted: string;
    export let ip: string;
    export let lastUsed: string;
    export let createdAt: string;

    function closeModal() {
        otpInputStatus = false;
        errorOtp = false;
        modal?.hide();
    }

    function betterTimeHandle(time: string) {
        if (time) {
            return betterTime(time);
        }
        return time;
    }

    async function logoutBtn() {
        if (!otpInputStatus) {
            otpInputStatus = true;
            return;
        }
        let {data, success} = await waitfetchData("http://127.0.0.1:8000/session/delete", "DELETE", {
            id: deviceId,
            otp: $otpCodeInputValue
        });
        if (!success) {
            errorOtp = true;
            return;
        }
        errorOtp = false;
        dispatch("updateDeviceList");
        closeModal();
    }
</script>

<div id={modalId} data-modal-target="{modalId}" tabindex="-1" aria-hidden="true" class="absolute hidden top-1/2 sm:top-0 -translate-y-1/4 sm:-translate-y-0 overflow-y-auto overflow-x-hidden right-0 left-0 z-50 justify-center items-center w-full md:inset-0 h-modal md:h-full">
    <div class="relative p-4 w-full max-w-xl h-full md:h-auto">
        <!-- Modal content -->
        <div class="relative p-4 bg-white rounded-lg shadow dark:bg-gray-800 sm:p-5">
                <!-- Modal header -->
                <div class="flex justify-between mb-4 rounded-t sm:mb-5">
                    <div class="text-lg text-gray-900 md:text-xl dark:text-white">
                        <p>
                            {userAgentFormatted}
                        </p>
                    </div>
                    <div>
                        <button type="button" class="text-gray-400 bg-transparent hover:bg-gray-200 hover:text-gray-900 rounded-lg text-sm p-1.5 inline-flex dark:hover:bg-gray-600 dark:hover:text-white" on:click={closeModal}>
                            <svg aria-hidden="true" class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg"><path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd"></path></svg>
                            <span class="sr-only">Close modal</span>
                        </button>
                    </div>
                </div>
                {#if !otpInputStatus}
                    <dl>
                        <dt class="mb-2 font-semibold leading-none text-gray-900 dark:text-white">Details</dt>
                        <dd class="mb-4 font-light text-gray-500 sm:mb-5 dark:text-gray-400">{userAgent}</dd>
                        <dt class="mb-2 font-semibold leading-none text-gray-900 dark:text-white">Ip address</dt>
                        <dd class="mb-4 font-light text-gray-500 sm:mb-5 dark:text-gray-400">{ip}</dd>
                        <dt class="mb-2 font-semibold leading-none text-gray-900 dark:text-white">Login date</dt>
                        <dd class="mb-4 font-light text-gray-500 sm:mb-5 dark:text-gray-400">{betterTimeHandle(createdAt)}</dd>
                        <dt class="mb-2 font-semibold leading-none text-gray-900 dark:text-white">Last action date</dt>
                        <dd class="mb-4 font-light text-gray-500 sm:mb-5 dark:text-gray-400">{betterTimeHandle(lastUsed)}</dd>
                    </dl>
                {:else}
                    <dl>
                        <OtpInput errorOtp={errorOtp} errorOtpText={"Otp code not valid"}/>
                    </dl>
                {/if}
                {#if !activeSession}
                    <div class="flex justify-end items-center">
                        <button on:click={logoutBtn} type="button" class="inline-flex items-center text-white bg-red-600 hover:bg-red-700 focus:ring-4 focus:outline-none focus:ring-red-300 font-medium rounded-lg text-sm px-5 py-2.5 text-center dark:bg-red-500 dark:hover:bg-red-600 dark:focus:ring-red-900">
                            <svg aria-hidden="true" class="w-5 h-5 mr-1.5 -ml-1" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg"><path fill-rule="evenodd" d="M9 2a1 1 0 00-.894.553L7.382 4H4a1 1 0 000 2v10a2 2 0 002 2h8a2 2 0 002-2V6a1 1 0 100-2h-3.382l-.724-1.447A1 1 0 0011 2H9zM7 8a1 1 0 012 0v6a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v6a1 1 0 102 0V8a1 1 0 00-1-1z" clip-rule="evenodd"></path></svg>
                            Logout
                        </button>
                    </div>
                {/if}
        </div>
    </div>
</div>