<script lang="ts">
	import { masterPassword } from '$lib/store/passwordStore';
    import { otpCodeInputValue } from '$lib/store/otpStore';
	import { getSalt, login } from '$lib/logic/login';
    import { checkMail } from '$lib/logic/utils';
	import OtpInput from '$lib/components/OtpInput.svelte';

    let loginStep: number = 1;
    
    // otp code
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

    async function loginBtn() {
        if ($otpCodeInputValue.length!= 6) {
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
        let {success, message} = await login(email, password, salt.toString(), $otpCodeInputValue);
        if (!success) {
            errorOtp = true;
            errorOtpText = message;
            return;
        }
        
        console.log('Login success');
        masterPassword.set(password);
        setTimeout(() => {
            window.location.href = "/passwords";
        }, 100);
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
                        <div class="mb-2 mx-auto">
                            <OtpInput errorOtp={errorOtp} errorOtpText={errorOtpText} />
                        </div>                        
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