<script lang="ts">
    import { otpCodeInputValue } from "$lib/store/otpStore";

    let code1: string;
	let code2: string;
	let code3: string;
	let code4: string;
	let code5: string;
	let code6: string;

	export let errorOtp: boolean;
	export let errorOtpText: string;

    function getOtpCode() {
		return `${code1}${code2}${code3}${code4}${code5}${code6}`;
	}

    function updateOtpCode () {
        console.log(getOtpCode());
        otpCodeInputValue.set(getOtpCode());
    }
</script>

<div class="mx-auto max-w-sm">
    <div class="mx-auto mb-2 flex space-x-2 rtl:space-x-reverse">
        <div>
            <label for="code-1" class="sr-only">First code</label>
            <input
                bind:value={code1}
                on:keyup={updateOtpCode}
                type="text"
                maxlength="1"
                data-focus-input-init
                data-focus-input-next="code-2"
                id="code-1"
                class="focus:ring-primary-500 focus:border-primary-500 dark:focus:ring-primary-500 dark:focus:border-primary-500 block h-10 w-10 rounded-lg border border-gray-300 bg-white py-3 text-center text-sm font-extrabold text-gray-900 md:h-14 md:w-14 dark:border-gray-600 dark:bg-gray-700 dark:text-white dark:placeholder-gray-400"
            />
        </div>
        <div>
            <label for="code-2" class="sr-only">Second code</label>
            <input
                bind:value={code2}
                on:keyup={updateOtpCode}
                type="text"
                maxlength="1"
                data-focus-input-init
                data-focus-input-prev="code-1"
                data-focus-input-next="code-3"
                id="code-2"
                class="focus:ring-primary-500 focus:border-primary-500 dark:focus:ring-primary-500 dark:focus:border-primary-500 block h-10 w-10 rounded-lg border border-gray-300 bg-white py-3 text-center text-sm font-extrabold text-gray-900 md:h-14 md:w-14 dark:border-gray-600 dark:bg-gray-700 dark:text-white dark:placeholder-gray-400"
            />
        </div>
        <div>
            <label for="code-3" class="sr-only">Third code</label>
            <input
                bind:value={code3}
                on:keyup={updateOtpCode}
                type="numer"
                maxlength="1"
                data-focus-input-init
                data-focus-input-prev="code-2"
                data-focus-input-next="code-4"
                id="code-3"
                class="focus:ring-primary-500 focus:border-primary-500 dark:focus:ring-primary-500 dark:focus:border-primary-500 block h-10 w-10 rounded-lg border border-gray-300 bg-white py-3 text-center text-sm font-extrabold text-gray-900 md:h-14 md:w-14 dark:border-gray-600 dark:bg-gray-700 dark:text-white dark:placeholder-gray-400"
            />
        </div>
        <div>
            <label for="code-4" class="sr-only">Fourth code</label>
            <input
                bind:value={code4}
                on:keyup={updateOtpCode}
                type="numer"
                maxlength="1"
                data-focus-input-init
                data-focus-input-prev="code-3"
                data-focus-input-next="code-5"
                id="code-4"
                class="focus:ring-primary-500 focus:border-primary-500 dark:focus:ring-primary-500 dark:focus:border-primary-500 block h-10 w-10 rounded-lg border border-gray-300 bg-white py-3 text-center text-sm font-extrabold text-gray-900 md:h-14 md:w-14 dark:border-gray-600 dark:bg-gray-700 dark:text-white dark:placeholder-gray-400"
            />
        </div>
        <div>
            <label for="code-5" class="sr-only">Fifth code</label>
            <input
                bind:value={code5}
                on:keyup={updateOtpCode}
                type="numer"
                maxlength="1"
                data-focus-input-init
                data-focus-input-prev="code-4"
                data-focus-input-next="code-6"
                id="code-5"
                class="focus:ring-primary-500 focus:border-primary-500 dark:focus:ring-primary-500 dark:focus:border-primary-500 block h-10 w-10 rounded-lg border border-gray-300 bg-white py-3 text-center text-sm font-extrabold text-gray-900 md:h-14 md:w-14 dark:border-gray-600 dark:bg-gray-700 dark:text-white dark:placeholder-gray-400"
            />
        </div>
        <div>
            <label for="code-6" class="sr-only">Sixth code</label>
            <input
                bind:value={code6}
                on:keyup={updateOtpCode}
                type="numer"
                maxlength="1"
                data-focus-input-init
                data-focus-input-prev="code-5"
                id="code-6"
                class="focus:ring-primary-500 focus:border-primary-500 dark:focus:ring-primary-500 dark:focus:border-primary-500 block h-10 w-10 rounded-lg border border-gray-300 bg-white py-3 text-center text-sm font-extrabold text-gray-900 md:h-14 md:w-14 dark:border-gray-600 dark:bg-gray-700 dark:text-white dark:placeholder-gray-400"
            />
        </div>
    </div>
    <p
        id="helper-text-explanation"
        class="mb-4 mt-2 text-sm text-gray-500 dark:text-gray-400"
    >
        Please introduce the 6 digit code from Autenticator App.
    </p>
    {#if errorOtp}
        <p class="mb-2 text-sm text-red-600 dark:text-red-500">{errorOtpText}</p>
    {/if}
</div>