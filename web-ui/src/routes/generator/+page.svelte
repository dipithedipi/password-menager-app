<script lang="ts">
	import { passwordStrength } from 'check-password-strength'

	let randomPassword: string = '';
	let passwordLength: number = 32;
	let lowercase: boolean = true;
	let uppercase: boolean = true;
	let numbers: boolean = true;
	let specialCharacters: boolean = true;
	let securityLevel: number = 1;
	let securityLevelMap: any = {
		0: ["bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-300" ,'To Weak'],
		1: ["bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-300" ,'Weak'],
		2: ["bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-300" ,'Medium'],
		3: ["bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300" ,'Strong'],
	};

	let options: {
		id: number;
		value: string;
		minDiversity: number;
		minLength: number;
	}[]

	function generateRandomPassword() {
		if (!lowercase && !uppercase && !numbers && !specialCharacters) {
			alert('Please select at least one option');
			return;
		};

		let characters = '';
		if (lowercase) characters += 'abcdefghijklmnopqrstuvwxyz';
		if (uppercase) characters += 'ABCDEFGHIJKLMNOPQRSTUVWXYZ';
		if (numbers) characters += '0123456789';
		if (specialCharacters) characters += '!@#$%^&*()_+{}:"<>?|[];\',./`~';
		if (characters === '') return;

		let password = '';
		for (let i = 0; i < passwordLength; i++) {
			password += characters.charAt(Math.floor(Math.random() * characters.length));
		}
		randomPassword = password;
		calculateStrenght();
	}

	function copyToClipboard() {
		navigator.clipboard.writeText(randomPassword);
	}

	function calculateStrenght() {
		securityLevel = passwordStrength(randomPassword).id;
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

<div class="h-screen justify-center bg-gray-700">
	<div class="pt-8 sm:pt-7 mx-2">
		<div
			class="mx-auto flex max-w-lg items-center rounded-lg border border-gray-700 bg-gray-800 p-3"
		>
			<div class="relative w-full">
				<div class="pointer-events-none absolute inset-y-0 start-0 flex items-center ps-3">
					<svg
						class="h-5 w-5 text-gray-500 dark:text-gray-400"
						aria-hidden="true"
						xmlns="http://www.w3.org/2000/svg"
						fill="currentColor"
						viewBox="0 0 24 24"
					>
						<path d="M10.313 11.5656L18.253 3.62561L20.3744 5.74693L18.9602 7.16114L21.0815 9.28246L17.5459 12.818L15.4246 10.6967L12.4343 13.687C13.4182 15.5719 13.1186 17.9524 11.5355 19.5355C9.58291 21.4881 6.41709 21.4881 4.46447 19.5355C2.51184 17.5829 2.51184 14.4171 4.46447 12.4644C6.04755 10.8814 8.42809 10.5818 10.313 11.5656ZM9.41421 17.4142C10.1953 16.6331 10.1953 15.3668 9.41421 14.5858C8.63316 13.8047 7.36684 13.8047 6.58579 14.5858C5.80474 15.3668 5.80474 16.6331 6.58579 17.4142C7.36684 18.1952 8.63316 18.1952 9.41421 17.4142Z"></path>
					</svg>
				</div>
				<input
					type="text"
					id="random-password"
					bind:value={randomPassword}
					class="block w-full rounded-lg border border-gray-300 bg-gray-50 p-2.5 ps-10 text-sm text-gray-900 focus:border-blue-500 focus:ring-blue-500 dark:border-gray-600 dark:bg-gray-700 dark:text-white dark:placeholder-gray-400 dark:focus:border-blue-500 dark:focus:ring-blue-500"
					placeholder="Random Password"
				/>
			</div>
			<button
				type="button"
				on:click={generateRandomPassword}
				class="ms-2 inline-flex items-center rounded-lg border border-blue-700 bg-blue-700 px-3 py-2.5 text-sm font-medium text-white hover:bg-blue-800 focus:outline-none focus:ring-4 focus:ring-blue-300 dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800"
			>
				<svg
					class="me-2 h-5 w-5"
					aria-hidden="true"
					xmlns="http://www.w3.org/2000/svg"
					fill="currentColor"
					viewBox="0 0 24 24"
				>
					<path d="M10.9979 1.58018C11.6178 1.22132 12.3822 1.22132 13.0021 1.58018L20.5021 5.92229C21.1197 6.27987 21.5 6.93946 21.5 7.65314V16.3469C21.5 17.0606 21.1197 17.7202 20.5021 18.0778L13.0021 22.4199C12.3822 22.7788 11.6178 22.7788 10.9979 22.4199L3.49793 18.0778C2.88029 17.7202 2.5 17.0606 2.5 16.3469V7.65314C2.5 6.93947 2.88029 6.27987 3.49793 5.92229L10.9979 1.58018ZM4.5 7.65314V7.65792L11.0021 11.4223C11.6197 11.7799 12 12.4395 12 13.1531V20.689L19.5 16.3469V7.65314L12 3.31104L4.5 7.65314ZM6.13208 12.3C6.13206 11.7477 5.74432 11.0761 5.26604 10.7999C4.78776 10.5238 4.40004 10.7476 4.40006 11.2999C4.40008 11.8522 4.78782 12.5238 5.2661 12.7999C5.74439 13.0761 6.1321 12.8523 6.13208 12.3ZM8.72899 18.7982C9.20728 19.0743 9.59499 18.8505 9.59497 18.2982C9.59495 17.7459 9.20721 17.0743 8.72893 16.7982C8.25065 16.522 7.86293 16.7459 7.86295 17.2982C7.86297 17.8504 8.25071 18.522 8.72899 18.7982ZM5.2661 16.799C5.74439 17.0751 6.1321 16.8513 6.13208 16.299C6.13206 15.7467 5.74432 15.0751 5.26604 14.799C4.78776 14.5228 4.40004 14.7467 4.40006 15.2989C4.40008 15.8512 4.78782 16.5228 5.2661 16.799ZM8.72851 14.7995C9.20679 15.0756 9.5945 14.8518 9.59448 14.2995C9.59446 13.7472 9.20673 13.0756 8.72844 12.7995C8.25016 12.5233 7.86245 12.7471 7.86246 13.2994C7.86248 13.8517 8.25022 14.5233 8.72851 14.7995ZM14.8979 8.00001C15.3762 7.72388 15.3762 7.27619 14.8979 7.00006C14.4196 6.72394 13.6441 6.72394 13.1658 7.00006C12.6875 7.27619 12.6875 7.72388 13.1658 8.00001C13.6441 8.27614 14.4196 8.27614 14.8979 8.00001ZM10.0981 7.00006C10.5764 7.27619 10.5764 7.72388 10.0981 8.00001C9.61982 8.27614 8.84434 8.27614 8.36604 8.00001C7.88774 7.72388 7.88774 7.27619 8.36604 7.00006C8.84434 6.72394 9.61982 6.72394 10.0981 7.00006ZM15.9954 15.3495C16.5932 15.0043 17.0779 14.1649 17.0779 13.4745C17.0779 12.7842 16.5933 12.5044 15.9955 12.8496C15.3977 13.1948 14.9131 14.0342 14.913 14.7246C14.913 15.4149 15.3976 15.6947 15.9954 15.3495Z"></path>
				</svg>New
			</button>
			<button
				type="button"
				on:click={copyToClipboard}
				class="ms-2 inline-flex items-center rounded-lg border border-blue-700 bg-blue-700 px-3 py-2.5 text-sm font-medium text-white hover:bg-blue-800 focus:outline-none focus:ring-4 focus:ring-blue-300 dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800"
			>
				<svg
					class="sm:me-2 h-5 w-5"
					aria-hidden="true"
					xmlns="http://www.w3.org/2000/svg"
					fill="currentColor"
					viewBox="0 0 24 24"
				>
					<path d="M7 4V2H17V4H20.0066C20.5552 4 21 4.44495 21 4.9934V21.0066C21 21.5552 20.5551 22 20.0066 22H3.9934C3.44476 22 3 21.5551 3 21.0066V4.9934C3 4.44476 3.44495 4 3.9934 4H7ZM7 6H5V20H19V6H17V8H7V6ZM9 4V6H15V4H9Z"
					></path>
				</svg>
				<div class="hidden sm:block">
					Copy
				</div>
			</button>
		</div>
	</div>

	<div class="mx-2 md:mx-auto max-w-lg pt-2">
		<div class="relative max-w-lg rounded-lg border-gray-700 bg-gray-800 p-3">
            <div
				class="flex items-center font-medium rounded-lg border border-gray-300 bg-gray-50 p-2.5 py-3 text-sm text-gray-900 focus:border-blue-500 focus:ring-blue-500 dark:border-gray-600 dark:bg-gray-600 dark:text-white dark:placeholder-gray-400 dark:focus:border-blue-500 dark:focus:ring-blue-500"
			>
				Lenght:

				<div class="absolute right-6 flex max-w-[8rem] items-center">
					<button
						type="button"
						on:click={() => (passwordLength > 0 ? passwordLength -= 1 : 0)}
						id="decrement-button"
						data-input-counter-decrement="quantity-input"
						class="h-10 rounded-s-lg border border-gray-300 bg-gray-100 p-2 hover:bg-gray-200 focus:outline-none focus:ring-2 focus:ring-gray-100 dark:border-gray-600 dark:bg-gray-700 dark:hover:bg-gray-600 dark:focus:ring-gray-700"
					>
						<svg
							class="h-3 w-3 text-gray-900 dark:text-white"
							aria-hidden="true"
							xmlns="http://www.w3.org/2000/svg"
							fill="none"
							viewBox="0 0 18 2"
						>
							<path
								stroke="currentColor"
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M1 1h16"
							/>
						</svg>
					</button>
					<input
						type="number"
						bind:value={passwordLength}
						id="quantity-input"
						data-input-counter
						aria-describedby="helper-text-explanation"
						class="block h-10 w-14 border-x-0 border-gray-300 bg-gray-50 py-2.5 text-center text-sm text-gray-900 focus:border-blue-500 focus:ring-blue-500 dark:border-gray-600 dark:bg-gray-700 dark:text-white dark:placeholder-gray-400 dark:focus:border-blue-500 dark:focus:ring-blue-500"
						placeholder="32"
					/>
					<button
						type="button"
						on:click={() => (passwordLength += 1)}
						id="increment-button"
						data-input-counter-increment="quantity-input"
						class="h-10 rounded-e-lg border border-gray-300 bg-gray-100 p-2 hover:bg-gray-200 focus:outline-none focus:ring-2 focus:ring-gray-100 dark:border-gray-600 dark:bg-gray-700 dark:hover:bg-gray-600 dark:focus:ring-gray-700"
					>
						<svg
							class="h-3 w-3 text-gray-900 dark:text-white"
							aria-hidden="true"
							xmlns="http://www.w3.org/2000/svg"
							fill="none"
							viewBox="0 0 18 18"
						>
							<path
								stroke="currentColor"
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M9 1v16M1 9h16"
							/>
						</svg>
					</button>
				</div>
			</div>

			
	<button id="dropdownMenuIconButton" data-dropdown-toggle="dropdownDots" class="w-full border text-sm border-gray-600 justify-center inline-flex items-center mt-3 p-1 font-normal text-center text-gray-900 bg-white rounded-lg hover:bg-gray-100 focus:ring-4 focus:outline-none dark:text-white focus:ring-gray-50 dark:bg-gray-800 dark:hover:bg-gray-700 dark:focus:ring-gray-600" type="button">
		Options
		<svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 mt-1 scale-110" viewBox="0 0 24 24" fill="currentColor">
			<path d="M12 15.0006L7.75732 10.758L9.17154 9.34375L12 12.1722L14.8284 9.34375L16.2426 10.758L12 15.0006Z"></path>
		</svg>
	</button>
	

	
	<!-- Dropdown menu -->
	<div id="dropdownDots" class="border border-gray-600 z-10 hidden rounded-lg w-44 dark:bg-gray-800 dark:divide-gray-600">
		<div class="mt-3 flex hover:bg-gray-700 items-center rounded-lg border border-gray-600 ps-4 mx-2">
			<input
				checked={lowercase}
				on:click={() => {lowercase = !lowercase}}
				id="bordered-checkbox-lowercase"
				type="checkbox"
				name="bordered-checkbox"
				class="h-4 w-4 rounded border-gray-300 bg-gray-100 text-blue-600 focus:ring-2 focus:ring-blue-500 dark:border-gray-600 dark:bg-gray-700 dark:ring-offset-gray-800 dark:focus:ring-blue-600"
			/>
			<label
				for="bordered-checkbox-lowercase"
				class="ms-2 w-full py-4 text-sm font-medium text-gray-900 dark:text-gray-300"
				>Lowercase</label
			>
		</div>
		<div class="mt-3 flex items-center hover:bg-gray-700 rounded-lg border border-gray-600 ps-4 mx-2">
			<input
				checked={uppercase}
				id="bordered-checkbox-uppercase"
				on:click={() => {uppercase = !uppercase}}
				type="checkbox"
				value=""
				name="bordered-checkbox"
				class="h-4 w-4 rounded border-gray-300 bg-gray-100 text-blue-600 focus:ring-2 focus:ring-blue-500 dark:border-gray-600 dark:bg-gray-700 dark:ring-offset-gray-800 dark:focus:ring-blue-600"
			/>
			<label
				for="bordered-checkbox-uppercase"
				class="ms-2 w-full py-4 text-sm font-medium text-gray-900 dark:text-gray-300"
				>Uppercase</label
			>
		</div>
		<div class="mt-3 flex items-center hover:bg-gray-700 rounded-lg border border-gray-600 ps-4 mx-2">
			<input
				checked={numbers}
				on:click={() => {numbers = !numbers}}
				id="bordered-checkbox-numbers"
				type="checkbox"
				name="bordered-checkbox"
				class="h-4 w-4 rounded border-gray-300 bg-gray-100 text-blue-600 focus:ring-2 focus:ring-blue-500 dark:border-gray-600 dark:bg-gray-700 dark:ring-offset-gray-800 dark:focus:ring-blue-600"
			/>
			<label
				for="bordered-checkbox-numbers"
				class="ms-2 w-full py-4 text-sm font-medium text-gray-900 dark:text-gray-300"
				>Numbers</label
			>
		</div>
		<div class="mt-3 mb-2 flex items-center hover:bg-gray-700 rounded-lg border border-gray-600 ps-4 mx-2">
			<input
				checked={specialCharacters}
				id="bordered-checkbox-special-characters"
				type="checkbox"
				on:click={() => {specialCharacters = !specialCharacters}}
				name="bordered-checkbox"
				class="h-4 w-4 rounded border-gray-300 bg-gray-100 text-blue-600 focus:ring-2 focus:ring-blue-500 dark:border-gray-600 dark:bg-gray-700 dark:ring-offset-gray-800 dark:focus:ring-blue-600"
			/>
			<label
				for="bordered-checkbox-special-characters"
				class="ms-2 w-full py-4 text-sm font-medium text-gray-900 dark:text-gray-300"
				>Special character</label
			>
		</div>
		
	</div> 
		</div>
        <div class="flex pt-3 items-center rounded-lg border-gray-700 bg-gray-800 p-3 mt-2"> 
            <div class="ms-2 w-full text-sm font-medium text-gray-900 dark:text-gray-300">
                Security level:
            </div>
            <span class="{securityLevelMap[securityLevel][0]} text-md font-medium me-2 px-2.5 py-0.5 rounded">{securityLevelMap[securityLevel][1]}</span>
        </div>
	</div>
</div>
