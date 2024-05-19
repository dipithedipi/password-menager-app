<script lang="ts">
	import { getCategory } from '$lib/logic/fetch';
  import { masterPassword } from '$lib/store/passwordStore';
  import { encryptAES } from '$lib/logic/cryptography';
  import { waitfetchData } from '$lib/logic/fetch';
	import { onMount, createEventDispatcher } from 'svelte';

  const dispatch = createEventDispatcher();

	export let modalId: string;

	let website: string = '';
	let details: string = '';
  let username: string = '';
	let password: string = '';
	let category: string = '';
  let otp: boolean = false;

  let errorUsername: boolean = false;
  let errorPassword: boolean = false;
  let errorWebsite: boolean = false;
  let errorCategory: boolean = false;
	let showPasswordValue: boolean = false;
	
  let categories: any[] = [];

	onMount(async () => {
		categories = await getCategory();
	});

	async function savePassword() {
    errorUsername = false;
    errorPassword = false;
    errorWebsite = false;
    errorCategory = false;

    if (website == '') {
      errorWebsite = true;
      return;
    }

    if (username == '') {
      errorUsername = true;
      return;
    }

    if (password == '') {
      errorPassword = true;
      return;
    }

    if (category == '') {
      errorCategory = true;
      return;
    }

    const encryptedPassword = encryptAES(password, $masterPassword);
    const encryptedDetails = encryptAES(details, $masterPassword);
    const encryptedUsername = encryptAES(username, $masterPassword);

    password = '';
    details = '';
    username = '';

    const body = {
      domain: website,
      username: encryptedUsername,
      description: encryptedDetails,
      password: encryptedPassword,
      category,
      otp
    };
    let {data, success} = await waitfetchData('http://127.0.0.1:8000/password/new', 'POST', body);
    if (!success) {
      alert("Error creating a new password");
      return;
    }
    console.log(data);
    cancel();
    // close the modal simulating click on cancel button
    console.log("dispach")
    dispatch('updatePasswords');
    document.getElementById("close-btn")?.click();
  }

	function showPassword() {
		showPasswordValue = !showPasswordValue;
	}

	function cancel() {
    website = '';
    username = '';
    details = '';
    password = '';
    category = '';
  }
</script>

<div
	id={modalId}
	data-modal-target={modalId}
	tabindex="-1"
	aria-hidden="true"
	class="h-modal absolute left-0 right-0 top-1/2 z-50 hidden w-full -translate-y-1/3 sm:translate-y-3 items-center justify-center overflow-y-auto overflow-x-hidden sm:top-0 md:inset-0 md:h-full"
>
	<div class="relative h-full w-full max-w-xl p-4 md:h-auto">
		<!-- Modal content -->
		<div class="relative rounded-lg bg-white p-4 shadow sm:p-5 dark:bg-gray-800">
			<!-- Modal header -->
			<div class="mb-4 flex justify-between rounded-t sm:mb-5">
				<div class="text-lg text-gray-900 md:text-xl dark:text-white">
					<h3 class="font-semibold">Save a new Password</h3>
				</div>
				<div>
					<button
            on:click={cancel}
						type="button"
						class="inline-flex rounded-lg bg-transparent p-1.5 text-sm text-gray-400 hover:bg-gray-200 hover:text-gray-900 dark:hover:bg-gray-600 dark:hover:text-white"
						data-modal-toggle={modalId}
					>
						<svg
							aria-hidden="true"
							class="h-5 w-5"
							fill="currentColor"
							viewBox="0 0 20 20"
							xmlns="http://www.w3.org/2000/svg"
							><path
								fill-rule="evenodd"
								d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
								clip-rule="evenodd"
							></path></svg
						>
						<span class="sr-only">Close modal</span>
					</button>
				</div>
			</div>
			<dl>
				<dt class="mb-2 font-semibold leading-none text-gray-900 dark:text-white">Website</dt>
				<dd class="mb-4 font-light text-gray-500 sm:mb-5 dark:text-gray-400">
					<div class="relative w-full">
						<div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
							<svg
								aria-hidden="true"
								class="h-5 w-5 scale-105 text-gray-500 dark:text-gray-400"
								fill="currentColor"
								viewbox="0 0 24 24"
								xmlns="http://www.w3.org/2000/svg"
							>
								<path
									d="M6.23509 6.45329C4.85101 7.89148 4 9.84636 4 12C4 16.4183 7.58172 20 12 20C13.0808 20 14.1116 19.7857 15.0521 19.3972C15.1671 18.6467 14.9148 17.9266 14.8116 17.6746C14.582 17.115 13.8241 16.1582 12.5589 14.8308C12.2212 14.4758 12.2429 14.2035 12.3636 13.3943L12.3775 13.3029C12.4595 12.7486 12.5971 12.4209 14.4622 12.1248C15.4097 11.9746 15.6589 12.3533 16.0043 12.8777C16.0425 12.9358 16.0807 12.9928 16.1198 13.0499C16.4479 13.5297 16.691 13.6394 17.0582 13.8064C17.2227 13.881 17.428 13.9751 17.7031 14.1314C18.3551 14.504 18.3551 14.9247 18.3551 15.8472V15.9518C18.3551 16.3434 18.3168 16.6872 18.2566 16.9859C19.3478 15.6185 20 13.8854 20 12C20 8.70089 18.003 5.8682 15.1519 4.64482C14.5987 5.01813 13.8398 5.54726 13.575 5.91C13.4396 6.09538 13.2482 7.04166 12.6257 7.11976C12.4626 7.14023 12.2438 7.12589 12.012 7.11097C11.3905 7.07058 10.5402 7.01606 10.268 7.75495C10.0952 8.2232 10.0648 9.49445 10.6239 10.1543C10.7134 10.2597 10.7307 10.4547 10.6699 10.6735C10.59 10.9608 10.4286 11.1356 10.3783 11.1717C10.2819 11.1163 10.0896 10.8931 9.95938 10.7412C9.64554 10.3765 9.25405 9.92233 8.74797 9.78176C8.56395 9.73083 8.36166 9.68867 8.16548 9.64736C7.6164 9.53227 6.99443 9.40134 6.84992 9.09302C6.74442 8.8672 6.74488 8.55621 6.74529 8.22764C6.74529 7.8112 6.74529 7.34029 6.54129 6.88256C6.46246 6.70541 6.35689 6.56446 6.23509 6.45329ZM12 22C6.47715 22 2 17.5228 2 12C2 6.47715 6.47715 2 12 2C17.5228 2 22 6.47715 22 12C22 17.5228 17.5228 22 12 22Z"
								></path>
							</svg>
						</div>
						<input
              				bind:value={website}
							type="text"
							class="focus:ring-primary-500 focus:border-primary-500 dark:focus:ring-primary-500 dark:focus:border-primary-500 block h-11 w-full rounded-lg border border-gray-300 bg-gray-50 p-2 pl-10 text-sm text-gray-900 dark:border-gray-600 dark:bg-gray-700 dark:text-white dark:placeholder-gray-400"
							placeholder="www.example.com"
						/>
					</div>
          {#if errorWebsite}
            <p class="text-sm font-semibold mt-1 text-red-600 dark:text-red-500">Website not valid</p>
          {/if}
				</dd>
        <dt class="mb-2 font-semibold leading-none text-gray-900 dark:text-white">Email or username</dt>
				<dd class="mb-4 font-light text-gray-500 sm:mb-5 dark:text-gray-400">
					<div class="relative w-full">
						<div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
							<svg
								aria-hidden="true"
								class="h-5 w-5 scale-105 text-gray-500 dark:text-gray-400"
								fill="currentColor"
								viewbox="0 0 24 24"
								xmlns="http://www.w3.org/2000/svg"
							>
                				<path d="M3 3H21C21.5523 3 22 3.44772 22 4V20C22 20.5523 21.5523 21 21 21H3C2.44772 21 2 20.5523 2 20V4C2 3.44772 2.44772 3 3 3ZM20 7.23792L12.0718 14.338L4 7.21594V19H20V7.23792ZM4.51146 5L12.0619 11.662L19.501 5H4.51146Z"></path>
							</svg>
						</div>
						<input
              				bind:value={username}
							type="text"
							class="focus:ring-primary-500 focus:border-primary-500 dark:focus:ring-primary-500 dark:focus:border-primary-500 block h-11 w-full rounded-lg border border-gray-300 bg-gray-50 p-2 pl-10 text-sm text-gray-900 dark:border-gray-600 dark:bg-gray-700 dark:text-white dark:placeholder-gray-400"
							placeholder="example@example.com"
						/>
					</div> 
          {#if errorUsername}
            <p class="text-sm font-semibold mt-1 text-red-600 dark:text-red-500">Username or Email not valid</p>
          {/if}
				</dd>
				<dt class="mb-2 font-semibold leading-none text-gray-900 dark:text-white">Details</dt>
				<dd class="mb-4 font-light text-gray-500 sm:mb-5 dark:text-gray-400">
					<div class="relative w-full">
						<div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
							<svg
								aria-hidden="true"
								class="h-5 w-5 scale-105 text-gray-500 dark:text-gray-400"
								fill="currentColor"
								viewbox="0 0 24 24"
								xmlns="http://www.w3.org/2000/svg"
							>
								<path
									d="M12 22C6.47715 22 2 17.5228 2 12C2 6.47715 6.47715 2 12 2C17.5228 2 22 6.47715 22 12C22 17.5228 17.5228 22 12 22ZM12 20C16.4183 20 20 16.4183 20 12C20 7.58172 16.4183 4 12 4C7.58172 4 4 7.58172 4 12C4 16.4183 7.58172 20 12 20ZM11 7H13V9H11V7ZM11 11H13V17H11V11Z"
								></path>
							</svg>
						</div>
						<input
              				bind:value={details}
							type="text"
							class="focus:ring-primary-500 focus:border-primary-500 dark:focus:ring-primary-500 dark:focus:border-primary-500 block h-11 w-full rounded-lg border border-gray-300 bg-gray-50 p-2 pl-10 text-sm text-gray-900 dark:border-gray-600 dark:bg-gray-700 dark:text-white dark:placeholder-gray-400"
							placeholder="Password info"
						/>
					</div>
				</dd>
				<dt class="mb-2 font-semibold leading-none text-gray-900 dark:text-white">Password</dt>
				<dd class="mb-4 font-light text-gray-500 sm:mb-5 dark:text-gray-400">
					<div class="relative w-full">
						<div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
							<svg
								aria-hidden="true"
								class="h-5 w-5 scale-105 text-gray-500 dark:text-gray-400"
								fill="currentColor"
								viewbox="0 0 24 24"
								xmlns="http://www.w3.org/2000/svg"
							>
								<path
									d="M17 14H12.6586C11.8349 16.3304 9.61244 18 7 18C3.68629 18 1 15.3137 1 12C1 8.68629 3.68629 6 7 6C9.61244 6 11.8349 7.66962 12.6586 10H23V14H21V18H17V14ZM7 14C8.10457 14 9 13.1046 9 12C9 10.8954 8.10457 10 7 10C5.89543 10 5 10.8954 5 12C5 13.1046 5.89543 14 7 14Z"
								></path>
							</svg>
						</div>
            {#if showPasswordValue == true}
              <input
                type="text"
                bind:value={password}
                class="focus:ring-primary-500 focus:border-primary-500 dark:focus:ring-primary-500 dark:focus:border-primary-500 block h-11 w-full rounded-lg border border-gray-300 bg-gray-50 p-2 pl-10 text-sm text-gray-900 dark:border-gray-600 dark:bg-gray-700 dark:text-white dark:placeholder-gray-400"
                placeholder="my&password"
              />
            {:else}
              <input
                type="password"
                bind:value={password}
                class="focus:ring-primary-500 focus:border-primary-500 dark:focus:ring-primary-500 dark:focus:border-primary-500 block h-11 w-full rounded-lg border border-gray-300 bg-gray-50 p-2 pl-10 text-sm text-gray-900 dark:border-gray-600 dark:bg-gray-700 dark:text-white dark:placeholder-gray-400"
                placeholder="•••••••••••"
                /> 
            {/if}
          </div>
          {#if errorPassword}
            <p class="text-sm font-semibold mt-1 text-red-600 dark:text-red-500">Password not valid</p>
          {/if}
				</dd>
				<dt class="mb-2 font-semibold leading-none text-gray-900 dark:text-white">Category</dt>
				<dd
					class="mb-2 flex space-x-3 font-light text-gray-500 sm:space-x-5 dark:text-gray-400"
				>
					<div class="relative w-80">
						<select
              				bind:value={category}
							id="countries"
							class="block h-11 w-full rounded-lg border border-gray-300 bg-gray-50 p-2.5 text-sm text-gray-900 focus:border-blue-500 focus:ring-blue-500 dark:border-gray-600 dark:bg-gray-700 dark:text-white dark:placeholder-gray-400 dark:focus:border-blue-500 dark:focus:ring-blue-500"
						>
							<option selected value="">Choose a category</option>
							{#each categories as category}
								<option value={category.name}>{category.name}</option>
							{/each}
						</select>
						{#if errorCategory}
						<p class="text-sm font-semibold mt-1 text-red-600 dark:text-red-500">Choose a category</p>
						{/if}
					</div>
					<div class="">
						<label class="mb-5 inline-flex h-11 cursor-pointer items-center">
							<input type="checkbox" bind:checked={otp} class="peer sr-only" />
							<div
								class="peer relative h-6 w-11 rounded-full bg-gray-200 after:absolute after:start-[2px] after:top-[2px] after:h-5 after:w-5 after:rounded-full after:border after:border-gray-300 after:bg-white after:transition-all after:content-[''] peer-checked:bg-blue-600 peer-checked:after:translate-x-full peer-checked:after:border-white peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 rtl:peer-checked:after:-translate-x-full dark:border-gray-600 dark:bg-gray-700 dark:peer-focus:ring-blue-800"
							></div>
							<span class="ms-3 text-sm font-bold text-gray-900 dark:text-gray-300">OTP secure</span
							>
						</label>
					</div>
				</dd>
			</dl>
			<div class="flex items-center justify-between">
				<div class="flex items-center space-x-2 sm:space-x-4">
					<button
						type="button"
            on:click={savePassword}
						class="bg-primary-700 hover:bg-primary-800 focus:ring-primary-300 dark:bg-primary-600 dark:hover:bg-primary-700 dark:focus:ring-primary-800 inline-flex items-center rounded-lg px-5 py-2.5 text-center text-sm font-medium text-white focus:outline-none focus:ring-4"
					>
						<svg
							aria-hidden="true"
							class="-ml-1 mr-1 h-5 w-5"
							fill="currentColor"
							viewBox="0 0 24 24"
							xmlns="http://www.w3.org/2000/svg"
						>
							<path
								d="M18 21V13H6V21H4C3.44772 21 3 20.5523 3 20V4C3 3.44772 3.44772 3 4 3H17L21 7V20C21 20.5523 20.5523 21 20 21H18ZM16 21H8V15H16V21Z"
							></path>
						</svg>
						Save
					</button>
					<button
						on:click={showPassword}
						type="button"
						class="hover:text-primary-700 rounded-lg border border-gray-200 bg-white px-5 py-2.5 text-sm font-medium text-gray-900 hover:bg-gray-100 focus:z-10 focus:outline-none focus:ring-4 focus:ring-gray-200 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-white dark:focus:ring-gray-700"
					>
						{#if showPasswordValue == true}
							Hide Password
						{:else}
							Show Password
						{/if}
					</button>
				</div>
				<button
					on:click={cancel}
					data-modal-toggle={modalId}
					type="button"
          id="close-btn"
					class="inline-flex items-center rounded-lg bg-red-600 px-5 py-2.5 text-center text-sm font-medium text-white hover:bg-red-700 focus:outline-none focus:ring-4 focus:ring-red-300 dark:bg-red-500 dark:hover:bg-red-600 dark:focus:ring-red-900"
				>
					<svg
						aria-hidden="true"
						class="-ml-1 mr-1.5 h-5 w-5"
						fill="currentColor"
						viewBox="0 0 20 20"
						xmlns="http://www.w3.org/2000/svg"
						><path
							fill-rule="evenodd"
							d="M9 2a1 1 0 00-.894.553L7.382 4H4a1 1 0 000 2v10a2 2 0 002 2h8a2 2 0 002-2V6a1 1 0 100-2h-3.382l-.724-1.447A1 1 0 0011 2H9zM7 8a1 1 0 012 0v6a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v6a1 1 0 102 0V8a1 1 0 00-1-1z"
							clip-rule="evenodd"
						></path></svg
					>
					Cancel
				</button>
			</div>
		</div>
	</div>
</div>
