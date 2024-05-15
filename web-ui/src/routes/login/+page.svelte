<script lang="ts">
	import { goto } from '$app/navigation';
	import { getSalt, login } from '$lib/logic/login';
    import { checkMail, setCookie } from '$lib/logic/utils';

    let loginStep: number = 1;
    
    // otp code
    let code1: string = '';
    let code2: string = '';
    let code3: string = '';
    let code4: string = '';
    let code5: string = '';
    let code6: string = '';
    let errorOtp: boolean = false;
    let errorOtpText: string = '';

    // email
    let email: string = '';
    let errorMail: boolean = false;
    let errorMailText: string = '';
    let emailRemember: boolean = false;
    if (typeof window!== 'undefined') {
        if (localStorage.getItem('emailRemember') == 'true') {
            email = localStorage.getItem('email') || '';
            emailRemember = true;
        }
    }

    // password
    let password: string = '';
    let errorPass: boolean = false;

    // salt
    let salt: string|boolean = '';

    async function nextStep() {
        // prevent default
        if (loginStep == 1) {
            // validate email and password
            if (!checkMail(email)) {
                errorMail = true;
                errorMailText = 'Email not valid';
                return;
            }

            if (password == '') {
                errorPass = true;
                return;
            }
        }

        // remember email
        if (emailRemember) {
            localStorage.setItem('email', email);
            localStorage.setItem('emailRemember', 'true');
        } else {
            localStorage.removeItem('email');
            localStorage.removeItem('emailRemember');
        }

        salt = await getSalt(email);
        if (salt === false) {
            // dispay user not found
            console.log('User not found');
            errorMail = true;
            errorMailText = 'User not found';
            return;
        }

        errorPass = false;
        errorMail = false;
        console.log(salt);
        loginStep++;
    }

    function getOtpCode() {
        return code1 + code2 + code3 + code4 + code5 + code6;
    }

    async function loginBtn() {
        const otpCode = getOtpCode();
        if (otpCode.length!= 6) {
            errorOtp = true;
            errorOtpText = "OTP code not complete";
            return;
        }

        // check if salt is not found
        if (salt === '' || salt === false) {
            alert("Salt not found, please try again.")
            return;
        }

        errorOtp = false;
        // Await the login function to resolve its promise
        let {success, message} = await login(email, password, salt.toString(), otpCode);
        if (!success) {
            errorOtp = true;
            errorOtpText = message;
            return;
        }
        
        console.log('Login success');
        // setTimeout(() => {
        //     window.location.href = "/passwords";
        // }, 1000);
    }
</script>

<section class="bg-gray-900">
    <div class="flex flex-col items-center h-screen justify-center px-6 py-8 mx-auto lg:py-0">
        <div class="w-full bg-white rounded-lg shadow dark:border md:mt-0 sm:max-w-md xl:p-0 dark:bg-gray-800 dark:border-gray-600">
            <div class="p-6 space-y-4 md:space-y-6 sm:p-8">
                <h1 class="text-xl font-bold leading-tight tracking-tight text-gray-900 md:text-2xl dark:text-white">
                    Sign in to your account
                </h1>
                {#if loginStep==1}
                    <div class="space-y-4 md:space-y-6">
                        <div>
                            <label for="email" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Your email</label>
                            <input bind:value={email} type="email" name="email" id="email" class="bg-gray-50 border border-gray-300 text-gray-900 sm:text-sm rounded-lg focus:ring-primary-600 focus:border-primary-600 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500" placeholder="name@company.com" >
                            {#if errorMail}
                                <p class="text-sm mt-1 text-red-600 dark:text-red-500">{errorMailText}</p>
                            {/if}
                        </div>
                        <div>
                            <label for="password" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Password</label>
                            <input bind:value={password} type="password" name="password" id="password" placeholder="••••••••" class="bg-gray-50 border border-gray-300 text-gray-900 sm:text-sm rounded-lg focus:ring-primary-600 focus:border-primary-600 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500" >
                            {#if errorPass}
                                <p class="text-sm mt-1 text-red-600 dark:text-red-500">Password not valid</p>
                            {/if}
                        </div>
                        <div class="flex items-center justify-between">
                            <div class="flex items-start">
                                <div class="flex items-center h-5">
                                <input bind:value={emailRemember} bind:checked={emailRemember} id="remember" aria-describedby="remember" type="checkbox" class="w-4 h-4 border border-gray-300 rounded bg-gray-50 focus:ring-3 focus:ring-primary-300 dark:bg-gray-700 dark:border-gray-600 dark:focus:ring-primary-600 dark:ring-offset-gray-800">
                                </div>
                                <div class="ml-3 text-sm">
                                <label for="remember" class="text-gray-500 dark:text-gray-300">Remember me</label>
                                </div>
                            </div>
                        </div>
                        <button on:click={nextStep} class="w-full text-white bg-primary-600 hover:bg-primary-700 focus:ring-4 focus:outline-none focus:ring-primary-300 font-medium rounded-lg text-sm px-5 py-2.5 text-center dark:bg-primary-600 dark:hover:bg-primary-700 dark:focus:ring-primary-800">Continue to sign in</button>
                    </div>
                {:else if loginStep==2}
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
                        <button on:click={loginBtn} class="w-full text-white bg-primary-600 hover:bg-primary-700 focus:ring-4 focus:outline-none focus:ring-primary-300 font-medium rounded-lg text-sm px-5 py-2.5 text-center dark:bg-primary-600 dark:hover:bg-primary-700 dark:focus:ring-primary-800">Sign in</button>
                    </div>
                {/if}
                <p class="text-sm font-light text-gray-500 dark:text-gray-400">
                    Don't have an account yet? <a href="/register" class="font-medium text-primary-600 hover:underline dark:text-primary-500">Sign up</a>
                </p>
            </div>
        </div>
    </div>
  </section>