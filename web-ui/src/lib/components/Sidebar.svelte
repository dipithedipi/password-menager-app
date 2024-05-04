<script>
    import { Sidebar, SidebarGroup, SidebarItem, SidebarWrapper, SidebarCta } from 'flowbite-svelte';
    let spanClass = 'flex-1 ms-3 whitespace-nowrap';
</script>

<div class="h-full">
    <Sidebar>
        <SidebarWrapper>
            <SidebarGroup>
                <SidebarItem label="Passwords">
                </SidebarItem>
                <SidebarItem label="Events" {spanClass}>
                </SidebarItem>
                <SidebarItem label="Devices" {spanClass}>
                </SidebarItem>
            </SidebarGroup>
            <SidebarGroup border>
                <SidebarItem label="Password Health">
                </SidebarItem>
                <SidebarItem label="Password Generator">
                </SidebarItem>
                <!-- bottom page -->
                <SidebarGroup border>
                    <div class="">
                        <SidebarCta label="Alpha">
                            <p class="mb-3 text-sm text-primary-900 dark:text-primary-400">This site is in development</p>
                        </SidebarCta>
                        <SidebarItem label="Settings"></SidebarItem>
                    </div>
                </SidebarGroup>
            </SidebarGroup>
        </SidebarWrapper>
    </Sidebar>
</div>


