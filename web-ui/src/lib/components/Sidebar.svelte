<script>
    import { Sidebar, SidebarGroup, SidebarItem, SidebarWrapper, SidebarCta } from 'flowbite-svelte';
    let spanClass = 'flex-1 ms-3 whitespace-nowrap';
</script>

<div>
    <Sidebar>
        <SidebarWrapper>
            <SidebarGroup>
                <SidebarItem label="Passwords" href="/passwords">
                </SidebarItem>
                <SidebarItem label="Events" {spanClass} href="/events">
                </SidebarItem>
                <SidebarItem label="Devices" {spanClass} href="/devices">
                </SidebarItem>
            </SidebarGroup>
            <SidebarGroup border>
                <SidebarItem label="Password Health" href="/health">
                </SidebarItem>
                <SidebarItem label="Password Generator" href="/generator">
                </SidebarItem>
                <!-- bottom page -->
                <SidebarGroup border>
                    <div>
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


