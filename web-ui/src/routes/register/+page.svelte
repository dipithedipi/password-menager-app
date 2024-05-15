<script lang="ts">
	import { checkMail, generateQRCode } from '$lib/logic/utils';
    import { checkUsername, register, verifyAccount } from '$lib/logic/register';

    let username = '';
    let errorUsername = false;
    let errorUsernameNotValid = false;

    let email = '';
    let errorEmail = false;

    let password = '';
    let errorPass = false;

    let confirmPassword = '';
    let errorConfirmPass = false;

    let termAndConditions = false;
    let errorTerms = false;

    let otpUrlQrCodeImage = ""

    let code1: string = '';
    let code2: string = '';
    let code3: string = '';
    let code4: string = '';
    let code5: string = '';
    let code6: string = '';
    let errorOtp: boolean = false;
    let errorOtpText: string = '';

    function getOtpCode() {
        return code1 + code2 + code3 + code4 + code5 + code6;
    }

    let registerStep = 0;
    async function nextStep() {
        errorEmail = false;
        errorPass = false;
        errorConfirmPass = false;
        errorTerms = false;
        errorUsernameNotValid = false;
        errorUsername = false;

        if (registerStep === 0) {
            if (username.length === 0) {
                errorUsernameNotValid = true;
                return;
            }

            const validUsername = await checkUsername(username);
            if (!validUsername) {
                errorUsername = true;
                return;
            }
            registerStep++;
        } else if (registerStep === 1) {
            if (!termAndConditions) {
                errorTerms = true;
                return;
            }
            if (!checkMail(email)) {
                errorEmail = true;
                return;
            }
            if (password.length < 0) {
                errorPass = true;
                return;
            }
            if (password !== confirmPassword) {
                errorConfirmPass = true;
                return;
            }

            // register
            let {success, message, otp} = await register(email, username, password);
            if (!success) {
                alert(message);
                return
            }

            // generate qrcode image
            otpUrlQrCodeImage = await generateQRCode("otpauth://totp/Password%2520Manager:alpha@chad.com?issuer=Password%2BManager&secret=FCOLLDRRKP2CCEIQPBVGPI3C52TPHGZQJFYB5QBP3VGYQ7RKHWLQ");
            console.log(otpUrlQrCodeImage);

            registerStep++;
        } else if (registerStep === 2) {
            // otp qr code
            registerStep++;
        } else if (registerStep === 3) {
            // otp code test
            let otpCode = getOtpCode();
            let {success, message} = await verifyAccount(email, otpCode);
            if (!success) {
                errorOtp = true;
                errorOtpText = message;
                return;
            }
            console.log("Account verified");
            // redirect to login after 1 seconds
            setTimeout(() => {
                window.location.href = "/login";
            }, 1000);
        }
    }
</script>

<section class="flex flex-col bg-gray-50 dark:bg-gray-900 h-screen items-center justify-center sm:pt-0">
    <div class="w-11/12 flex flex-col items-center justify-center px-2 py-8 mx-auto md:h-screen lg:py-0">
        <div class="w-full bg-white rounded-lg shadow dark:border md:mt-0 sm:max-w-md xl:p-0 dark:bg-gray-800 dark:border-gray-700">
            <div class="p-6 space-y-4 md:space-y-6 sm:p-8">
                <h1 class="text-xl font-bold leading-tight tracking-tight text-gray-900 md:text-2xl dark:text-white">
                    Create and account
                </h1>
                    <div class="space-y-4 md:space-y-6">
                        {#if registerStep === 0}
                        <div>
                            <div>
                                <label for="username" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Choose a username</label>
                                <input bind:value={username} type="text" name="username" id="username" class="bg-gray-50 border border-gray-300 text-gray-900 sm:text-sm rounded-lg focus:ring-primary-600 focus:border-primary-600 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500" placeholder="Username">
                                {#if errorUsername}
                                    <p class="text-sm mt-1 text-red-600 dark:text-red-500">Username alredy taken</p>
                                {/if}
                                {#if errorUsernameNotValid}
                                    <p class="text-sm mt-1 text-red-600 dark:text-red-500">Username not valid</p>
                                {/if}
                            </div>
                            <button on:click={nextStep} class="mt-2 w-full text-white bg-primary-600 hover:bg-primary-700 focus:ring-4 focus:outline-none focus:ring-primary-300 font-medium rounded-lg text-sm px-5 py-2.5 text-center dark:bg-primary-600 dark:hover:bg-primary-700 dark:focus:ring-primary-800">Continue</button>
                        </div>
                        {:else if registerStep === 1}
                            <div>
                                <div>
                                    <label for="email" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Your email</label>
                                    <input bind:value={email} type="email" name="email" id="email" class="bg-gray-50 border border-gray-300 text-gray-900 sm:text-sm rounded-lg focus:ring-primary-600 focus:border-primary-600 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500" placeholder="name@company.com">
                                    {#if errorEmail}
                                        <p class="text-sm mt-1 text-red-600 dark:text-red-500">Email not valid</p>
                                    {/if}
                                </div>
                                <div class="py-4">
                                    <label for="password" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Password</label>
                                    <input bind:value={password} type="password" name="password" id="password" placeholder="••••••••" class="bg-gray-50 border border-gray-300 text-gray-900 sm:text-sm rounded-lg focus:ring-primary-600 focus:border-primary-600 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500">
                                    {#if errorPass}
                                        <p class="text-sm mt-1 text-red-600 dark:text-red-500">Password must be at least 12 characters</p>
                                    {/if}
                                </div>
                                <div>
                                    <label for="confirm-password" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Confirm password</label>
                                    <input bind:value={confirmPassword} type="confirm-password" name="confirm-password" id="confirm-password" placeholder="••••••••" class="bg-gray-50 border border-gray-300 text-gray-900 sm:text-sm rounded-lg focus:ring-primary-600 focus:border-primary-600 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500">
                                    {#if errorConfirmPass}
                                        <p class="text-sm mt-1 text-red-600 dark:text-red-500">Passwords doesn't mach</p>
                                    {/if}
                                </div>
                                <div class="flex items-start pt-3 pb-1">
                                    <div class="flex items-center h-5">
                                        <input bind:value={termAndConditions} id="terms" aria-describedby="terms" type="checkbox" class="w-4 h-4 border border-gray-300 rounded bg-gray-50 focus:ring-3 focus:ring-primary-300 dark:bg-gray-700 dark:border-gray-600 dark:focus:ring-primary-600 dark:ring-offset-gray-800">
                                    </div>
                                    <div class="ml-3 text-sm">
                                    <label for="terms" class="font-light text-gray-500 dark:text-gray-300">I accept the <a class="font-medium text-primary-600 hover:underline dark:text-primary-500" href="#">Terms and Conditions</a></label>
                                    </div>
                                </div>
                                {#if errorTerms}
                                    <p class="text-sm mb-2 text-red-600 dark:text-red-500">You must accept the terms and conditions</p>
                                {/if}
                                <button on:click={nextStep} class="w-full text-white bg-primary-600 hover:bg-primary-700 focus:ring-4 focus:outline-none focus:ring-primary-300 font-medium rounded-lg text-sm px-5 py-2.5 text-center dark:bg-primary-600 dark:hover:bg-primary-700 dark:focus:ring-primary-800">Continue</button>
                            </div>
                        {:else if registerStep === 2}
                            <div>
                                <div class="font-light text-gray-500">
                                    QR CODE OTP
                                    <img class="w-10/12 mx-auto py-1" src="{otpUrlQrCodeImage}" alt="">
                                </div>
                                <div class="space-y-3 pt-2">
                                    <button on:click={nextStep} class="w-full text-white bg-primary-600 hover:bg-primary-700 focus:ring-4 focus:outline-none focus:ring-primary-300 font-medium rounded-lg text-sm px-5 py-2.5 text-center dark:bg-primary-600 dark:hover:bg-primary-700 dark:focus:ring-primary-800">Continue</button>
                                </div>
                            </div>
                        {:else if registerStep === 3}
                        <div class="max-w-sm mx-auto">
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
                            <button on:click={nextStep} class="w-full text-white bg-primary-600 hover:bg-primary-700 focus:ring-4 focus:outline-none focus:ring-primary-300 font-medium rounded-lg text-sm px-5 py-2.5 text-center dark:bg-primary-600 dark:hover:bg-primary-700 dark:focus:ring-primary-800">Continue</button>
                        </div>
                        {/if}
                        <p class="text-sm font-light text-gray-500 dark:text-gray-400">
                            Already have an account? <a href="/login" class="font-medium text-primary-600 hover:underline dark:text-primary-500">Login here</a>
                        </p>
                    </div>
            </div>
        </div>
    </div>
  </section>