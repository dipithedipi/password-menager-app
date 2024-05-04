<script>
    import { Navbar, NavBrand, NavLi, NavUl, NavHamburger, Avatar, Dropdown, DropdownItem, DropdownHeader, DropdownDivider } from 'flowbite-svelte';
  </script>
  
<Navbar rounded color="form">
    <NavBrand href="/">
      <img src="/logo.jpg" class="me-3 h-8 sm:h-9 rounded-full" alt="Flowbite Logo" />
      <span class="self-center whitespace-nowrap text-xl font-semibold dark:text-white">PassGuard</span>
    </NavBrand>
    <div class="flex items-center md:order-2">
      <Avatar id="avatar-menu" class="ri-account-circle-line" />
      <NavHamburger class1="w-full md:flex md:w-auto md:order-1" />
    </div>
    <Dropdown placement="bottom" triggeredBy="#avatar-menu">
      <DropdownHeader>
        <span class="block text-sm">USERNAME</span>
        <span class="block truncate text-sm font-medium">EMAIL</span>
      </DropdownHeader>
      <DropdownItem>Settings</DropdownItem>
      <DropdownDivider />
      <DropdownItem>Sign out</DropdownItem>
    </Dropdown>
    
      <!-- <NavUl>
        <NavLi href="/">Password</NavLi>
        <NavLi href="/about">Events</NavLi>
        <NavLi href="/docs/components/navbar">Devices</NavLi>
        <NavLi href="/pricing">Password Health</NavLi>
        <NavLi href="/contact">Password Generator</NavLi>
      </NavUl> -->
</Navbar>