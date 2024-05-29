<script lang="ts">
    import { otpCodeInputValue } from "$lib/store/otpStore";

    let otpCode: string = "";
	export let errorOtp: boolean;
	export let errorOtpText: string;


    function updateOtpCode () {
        otpCodeInputValue.set(otpCode.toString());
    }
</script>

<style>
    /* Chrome, Safari, Edge, Opera */
	input::-webkit-outer-spin-button,
	input::-webkit-inner-spin-button {
		-webkit-appearance: none;
		margin: 0;
	}

	/* Firefox */
	input[type=number] {
		-moz-appearance: textfield;
	}
</style>

<div class="mx-auto max-w-sm">
    <div class="mx-auto mb-2">
        <div class="mb-2  mx-auto">
            <label for="number-input" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Otp code</label>
            <input type="number" on:keyup={updateOtpCode} bind:value={otpCode} id="number-input" aria-describedby="helper-text-explanation" class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500" placeholder="00000"/>
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