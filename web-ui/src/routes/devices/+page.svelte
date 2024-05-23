<script lang="ts">
	import { formatUserAgent } from '$lib/logic/utils';
    import { Modal } from 'flowbite';
	import MyModalDevice from './../../lib/components/MyModalDevice.svelte';
	import { waitfetchData } from '$lib/logic/fetch';
    import { onMount } from 'svelte';
	import PhoneDevice from '$lib/components/PhoneDevice.svelte';

    let devices:any = [];
    let modal: any = null

    let CurrentUser: boolean;
    let LastUse: string;
    let UserAgent: string;
    let deviceId: string;
    let createdAt: string;
    let deviceIp: string;
    let userAgentFormatted: string;

    async function getEvents() {
        let {data, success} = await waitfetchData("http://127.0.0.1:8000/session/gets", "GET", {});
        if (!success) {
            console.log("Error fetching data")
            return;
        }
        devices = data.sessions;

        // sort by last used with active session first
        devices.sort((a: any, b: any) => {
            if (a.CurrentUser && !b.CurrentUser) {
                return -1;
            } else if (!a.CurrentUser && b.CurrentUser) {
                return 1;
            } else {
                return new Date(b.LastUse).getTime() - new Date(a.LastUse).getTime();
            }
        });
    }

    function openModal(event: any, index: number) {
        CurrentUser = devices[index].CurrentUser;
        LastUse = devices[index].LastUse;
        UserAgent = devices[index].UserAgent;
        userAgentFormatted = formatUserAgent(UserAgent);
        deviceId = devices[index].DatabaseElemID;
        createdAt = devices[index].CreatedAt;
        deviceIp = devices[index].IpAddress;
        modal?.show();
    }

    onMount(async () => {
        await getEvents();
        modal = new Modal(document.getElementById('device'));
    });
</script>

<div class="h-screen pt-6 bg-gray-700">
    <MyModalDevice on:updateDeviceList={getEvents} modalId="device" modal={modal} ip={deviceIp} userAgentFormatted={userAgentFormatted} activeSession={CurrentUser} lastUsed={LastUse} userAgent={UserAgent} deviceId={deviceId} createdAt={createdAt}></MyModalDevice>
    {#each devices as device, index}
        <PhoneDevice on:openDeviceModal={(e)=> {openModal(e, index)}} activeSession={device.CurrentUser} lastUsed={device.LastUse} userAgent={device.UserAgent}></PhoneDevice>
    {/each}
</div>