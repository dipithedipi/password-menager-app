<script lang="ts">
	import { otpCodeInputValue } from '$lib/store/otpStore';
	import { getCategory, waitfetchData } from '$lib/logic/fetch';
	import { Modal } from 'flowbite';
	import { createEventDispatcher, onMount } from 'svelte';
	import OtpInput from './OtpInput.svelte';
	import { masterPassword } from '$lib/store/passwordStore';
	import { decryptAES, encryptAES } from '$lib/logic/cryptography';

	const dispatchUpdate = createEventDispatcher();

	export let modalId: string;
	export let passwordId: string;
	export let username: string;
	export let title: string;
	export let description: string;
	export let category: string;
	export let otpProtected: boolean;
	export let otpCodeInput: boolean;
	export let categories: any[];

	let editMode: boolean = false;
  	let copyPasswordMode: boolean = false;
  	let deleteMode: boolean = false;

	let newUsername: string;
	let newPassword: string;
	let newDescription: string;
	let newCategory: string;
	let newOtpProtected: boolean;

	let errorOtp: boolean = false;
	let errorOtpText: string = '';

	let modalElem: HTMLDivElement | null = null;
	let modal: Modal | null = null;
	$: if (modalElem !== null) {
		modal = new Modal(modalElem);
	}

	$: if (username) {
		otpCodeInput = false;
	}

	const dispatch = createEventDispatcher<{ modalDetect: Modal }>();

	$: if (modal !== null) {
		otpCodeInput = false;
		dispatch('modalDetect', modal);
	}

	async function deletePassword() {
    copyPasswordMode = false;
    editMode = false;
    deleteMode = true;

		if (otpProtected && !otpCodeInput) {
			otpCodeInput = true;
			console.log('Delete password otp protected');
			return;
		}

		console.log('Delete password');
		let { data, success } = await waitfetchData('http://127.0.0.1:8000/password/delete', 'DELETE', {
			passwordId,
			otp: $otpCodeInputValue == '' ? '000000' : $otpCodeInputValue
		});
		if (success) {
			console.log('Password deleted');
			dispatchUpdate('updatePasswords');
      resetModal();
			modal?.hide();
		} else {
			console.log('Error deleting password');
			if (otpCodeInput) {
				errorOtp = true;
				errorOtpText = 'Invalid OTP code';
			}
		}
	}

	async function copyPassword() {
    deleteMode = false;
    editMode = false;
    copyPasswordMode = true;

		if (otpProtected && !otpCodeInput) {  
      errorOtp = false;
			otpCodeInput = true;
			console.log('Copy password otp protected');
			return;
		}

		console.log('Copy password');
		let { data, success } = await waitfetchData('http://127.0.0.1:8000/password/get', 'POST', {
			passwordId,
			domain: title,
			otp: $otpCodeInputValue == '' ? '000000' : $otpCodeInputValue
		});
		if (!success) {
			console.log('Error getting password');
			if (otpCodeInput) {
				errorOtp = true;
				errorOtpText = 'Invalid OTP code';
			}
			return;
		}
		navigator.clipboard.writeText(decryptAES(data.password[0].password, $masterPassword));
		resetModal();
    modal?.hide();
    console.log('Password copied');
	}

	async function editPassword() {
    copyPasswordMode = false;
    deleteMode = false;

    if (deleteMode || copyPasswordMode) {
      otpCodeInput = false;
    }

    if (editMode && otpProtected && !otpCodeInput) {
      errorOtp = false;
      otpCodeInput = true;
      console.log('Edit password otp protected');
      return;
    } else {
      if (!editMode) {
        categories = await getCategory();
        newUsername = username;
        newDescription = description;
        newCategory = category;
        newOtpProtected = otpProtected;
        newPassword = '';
        editMode = true;
        return;
      }
    }

    if (editMode && !otpCodeInput && otpProtected) {
      console.log("block otp")
      return
    }
    
    if (newUsername == '' || newDescription == '' || newCategory == '' || newPassword == '') {
      alert('Please fill all the fields');
      return;
    }

		console.log('Edit password');
    let { data, success } = await waitfetchData('http://127.0.0.1:8000/password/update', 'PUT', {
      passwordId,
      newUsername: encryptAES(newUsername, $masterPassword),
      newDescription: encryptAES(newDescription, $masterPassword),
      newCategory,
      otpProtected: newOtpProtected,
      newPassword: encryptAES(newPassword, $masterPassword),
      otp: $otpCodeInputValue == '' ? '000000' : $otpCodeInputValue
    });
    if (!success) {
      console.log('Error editing password');
      if (otpCodeInput) {
        errorOtp = true;
        errorOtpText = 'Invalid OTP code';
      }
      return;
    }
    console.log('Password edited');
    dispatchUpdate('updatePasswords');
    resetModal();
    modal?.hide();
	}

  function resetModal() {
    copyPasswordMode = false;
    otpCodeInput = false;
    deleteMode = false;
    editMode = false;
    modal?.hide();
  }
</script>

<div
	bind:this={modalElem}
	id={modalId}
	data-modal-target={modalId}
	tabindex="-1"
	aria-hidden="true"
	class="h-modal fixed left-0 right-0 top-1/2 z-50 hidden w-full -translate-y-1/4 items-center justify-center overflow-y-auto overflow-x-hidden sm:top-0 sm:-translate-y-0 md:inset-0 md:h-full"
>
	<div class="relative h-full w-full max-w-xl p-4 md:h-auto">
		<!-- Modal content -->
		<div class="relative rounded-lg bg-white p-4 shadow sm:p-5 dark:bg-gray-800">
			<!-- Modal header -->
			<div class="mb-4 flex justify-between rounded-t sm:mb-5">
				<div class="text-lg text-gray-900 md:text-xl dark:text-white">
					<h3 class="font-semibold h-11 py-1">
						{#if editMode}
              Edit: {title}
            {:else if copyPasswordMode}
                Copy: {title}
            {:else if deleteMode}
              Delete: {title}
            {:else}
              {title}
            {/if}
					</h3>
				</div>
				<div>
					<button
						type="button"
						class="inline-flex rounded-lg bg-transparent p-1.5 text-sm text-gray-400 hover:bg-gray-200 hover:text-gray-900 dark:hover:bg-gray-600 dark:hover:text-white"
						on:click={resetModal}
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
			{#if !otpCodeInput}
				<dl class={editMode ? "space-y-3":""}>
					<dt class="mb-2 font-semibold leading-none text-gray-900 dark:text-white">
						{username.includes('@') ? 'Email' : 'Username'}
					</dt>
					{#if !editMode}
						<dd class="mb-4 font-light text-gray-500 sm:mb-5 dark:text-gray-400">{username}</dd>
					{:else}
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
										d="M3 3H21C21.5523 3 22 3.44772 22 4V20C22 20.5523 21.5523 21 21 21H3C2.44772 21 2 20.5523 2 20V4C2 3.44772 2.44772 3 3 3ZM20 7.23792L12.0718 14.338L4 7.21594V19H20V7.23792ZM4.51146 5L12.0619 11.662L19.501 5H4.51146Z"
									></path>
								</svg>
							</div>
							<input
								bind:value={newUsername}
								type="text"
								class="focus:ring-primary-500 focus:border-primary-500 dark:focus:ring-primary-500 dark:focus:border-primary-500 block h-11 w-full rounded-lg border border-gray-300 bg-gray-50 p-2 pl-10 text-sm text-gray-900 dark:border-gray-600 dark:bg-gray-700 dark:text-white dark:placeholder-gray-400"
								placeholder="example@example.com"
							/>
						</div>
					{/if}
					<dt class="mb-2 font-semibold leading-none text-gray-900 dark:text-white">Details</dt>
					{#if !editMode}
						<dd class="mb-4 font-light text-gray-500 sm:mb-5 dark:text-gray-400">{description}</dd>
					{:else}
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
								bind:value={newDescription}
								type="text"
								class="focus:ring-primary-500 focus:border-primary-500 dark:focus:ring-primary-500 dark:focus:border-primary-500 block h-11 w-full rounded-lg border border-gray-300 bg-gray-50 p-2 pl-10 text-sm text-gray-900 dark:border-gray-600 dark:bg-gray-700 dark:text-white dark:placeholder-gray-400"
								placeholder="Password info"
							/>
						</div>
					{/if}
          {#if editMode}
          <dt class="mb-2 font-semibold leading-none text-gray-900 dark:text-white">Password</dt>
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
            <input
              type="password"
              bind:value={newPassword}
              class="focus:ring-primary-500 focus:border-primary-500 dark:focus:ring-primary-500 dark:focus:border-primary-500 block h-11 w-full rounded-lg border border-gray-300 bg-gray-50 p-2 pl-10 text-sm text-gray-900 dark:border-gray-600 dark:bg-gray-700 dark:text-white dark:placeholder-gray-400"
              placeholder="•••••••••••"
            />
          </div>
          {/if}
					<dt class="mb-2 font-semibold leading-none text-gray-900 dark:text-white">Category</dt>
					{#if !editMode}
						<dd class="mb-4 font-light text-gray-500 sm:mb-5 dark:text-gray-400">{category}</dd>
					{:else}
						<div class="flex space-x-3 font-light text-gray-500 sm:space-x-5 dark:text-gray-400">
							<div class="relative w-80">
								<select
									bind:value={newCategory}
									id="categories"
									class="block h-11 w-full rounded-lg border border-gray-300 bg-gray-50 p-2.5 text-sm text-gray-900 focus:border-blue-500 focus:ring-blue-500 dark:border-gray-600 dark:bg-gray-700 dark:text-white dark:placeholder-gray-400 dark:focus:border-blue-500 dark:focus:ring-blue-500"
								>
									<option selected value="">Choose a category</option>
									{#each categories as category}
										<option value={category.name}>{category.name}</option>
									{/each}
								</select>
							</div>
							<div class="">
								<label class="mb-5 inline-flex h-11 cursor-pointer items-center">
									<input type="checkbox" bind:checked={newOtpProtected} class="peer sr-only" />
									<div
										class="peer relative h-6 w-11 rounded-full bg-gray-200 after:absolute after:start-[2px] after:top-[2px] after:h-5 after:w-5 after:rounded-full after:border after:border-gray-300 after:bg-white after:transition-all after:content-[''] peer-checked:bg-blue-600 peer-checked:after:translate-x-full peer-checked:after:border-white peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 rtl:peer-checked:after:-translate-x-full dark:border-gray-600 dark:bg-gray-700 dark:peer-focus:ring-blue-800"
									></div>
									<span class="ms-3 text-sm font-bold text-gray-900 dark:text-gray-300"
										>OTP secure</span
									>
								</label>
							</div>
						</div>
					{/if}
				</dl>
			{:else}
				<dl>
					<OtpInput {errorOtp} {errorOtpText} />
				</dl>
			{/if}
			<div class="flex items-center justify-between">
				<div class="flex items-center space-x-5 sm:space-x-4">
					<button
						on:click={editPassword}
						type="button"
						class="bg-primary-700 hover:bg-primary-800 focus:ring-primary-300 dark:bg-primary-600 dark:hover:bg-primary-700 dark:focus:ring-primary-800 inline-flex items-center rounded-lg px-5 py-2.5 text-center text-sm font-medium text-white focus:outline-none focus:ring-4"
					>
						<svg
							aria-hidden="true"
							class="-ml-1 mr-1 h-5 w-5"
							fill="currentColor"
							viewBox="0 0 20 20"
							xmlns="http://www.w3.org/2000/svg"
							><path
								d="M17.414 2.586a2 2 0 00-2.828 0L7 10.172V13h2.828l7.586-7.586a2 2 0 000-2.828z"
							></path><path
								fill-rule="evenodd"
								d="M2 6a2 2 0 012-2h4a1 1 0 010 2H4v10h10v-4a1 1 0 112 0v4a2 2 0 01-2 2H4a2 2 0 01-2-2V6z"
								clip-rule="evenodd"
							></path></svg
						>
            {editMode ? 'Save' : 'Edit'}
					</button>
					<button
						on:click={copyPassword}
						type="button"
						class="hover:text-primary-700 inline-flex rounded-lg border border-gray-200 bg-white px-5 py-2.5 text-sm font-medium text-gray-900 hover:bg-gray-100 focus:z-10 focus:outline-none focus:ring-4 focus:ring-gray-200 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-white dark:focus:ring-gray-700"
					>
						<svg
							xmlns="http://www.w3.org/2000/svg"
							class="-ml-1 mr-1 h-5 w-5"
							viewBox="0 0 24 24"
							fill="currentColor"
							><path
								d="M7 4V2H17V4H20.0066C20.5552 4 21 4.44495 21 4.9934V21.0066C21 21.5552 20.5551 22 20.0066 22H3.9934C3.44476 22 3 21.5551 3 21.0066V4.9934C3 4.44476 3.44495 4 3.9934 4H7ZM7 6H5V20H19V6H17V8H7V6ZM9 4V6H15V4H9Z"
							></path></svg
						>
						Copy
						<!-- space -->
						<div class="ml-1 hidden sm:block">the password</div>
					</button>
				</div>
				<button
					type="button"
					on:click={deletePassword}
					class="inline-flex items-center rounded-lg bg-red-600 px-5 py-2.5 text-center text-sm font-medium text-white hover:bg-red-700 focus:outline-none focus:ring-4 focus:ring-red-300 dark:bg-red-500 dark:hover:bg-red-600 dark:focus:ring-red-900"
				>
					<svg
						aria-hidden="true"
						class="h-5 w-5 sm:-ml-1 sm:mr-1.5"
						fill="currentColor"
						viewBox="0 0 20 20"
						xmlns="http://www.w3.org/2000/svg"
						><path
							fill-rule="evenodd"
							d="M9 2a1 1 0 00-.894.553L7.382 4H4a1 1 0 000 2v10a2 2 0 002 2h8a2 2 0 002-2V6a1 1 0 100-2h-3.382l-.724-1.447A1 1 0 0011 2H9zM7 8a1 1 0 012 0v6a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v6a1 1 0 102 0V8a1 1 0 00-1-1z"
							clip-rule="evenodd"
						></path></svg
					>
					<div>
            {!deleteMode ? 'Delete' : 'Confirm'}
          </div>
				</button>
			</div>
		</div>
	</div>
</div>
