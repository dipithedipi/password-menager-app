import QRCode from 'qrcode'
import { parseISO, format } from 'date-fns';

function checkMail(mail: string): boolean {
  return /^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,6}$/.test(mail);
}

function setCookie(data: string) {
  document.cookie = data;
}

function calculateTimeDifference(lastUse: string): string {
  // Parse the lastUse string to a Date object
  const lastUseDate = parseISO(lastUse);

  // Get the current date and time
  const now = new Date();

  // Calculate the difference in hours
  const monthsDifference = Math.floor(now.getMonth() - lastUseDate.getMonth());
  const daysDifference = Math.floor(now.getDate() - lastUseDate.getDate());
  const hoursDifference = Math.floor(now.getHours() - lastUseDate.getHours());
  const minutesDifference = Math.floor(now.getMinutes() - lastUseDate.getMinutes());

  // Return the difference in a human readable format
  if (monthsDifference > 0) {
    return `${monthsDifference} months ago`;
  } else if (daysDifference > 0) {
    return `${daysDifference} days ago`;
  } else if (hoursDifference > 0) {
    return `${hoursDifference} hours ago`;
  } else {
    return `${minutesDifference} minutes ago`;
  }
}

function getCookie(cname: string) {
  let name = cname + "=";
  let decodedCookie = decodeURIComponent(document.cookie);
  let ca = decodedCookie.split(';');
  for(let i = 0; i <ca.length; i++) {
    let c = ca[i];
    while (c.charAt(0) == ' ') {
      c = c.substring(1);
    }
    if (c.indexOf(name) == 0) {
      return c.substring(name.length, c.length);
    }
  }
  return "";
}

async function generateQRCode(data: string): Promise<string> {
  return await QRCode.toDataURL(data);
}

function formatDate(date: string): string {
  //dd/MM/yyyy -> yyyy-MM-dd
  const dateArray = date.split("/");
  return `${dateArray[2]}-${dateArray[1]}-${dateArray[0]}`;
}

export { checkMail, formatDate, setCookie, generateQRCode, getCookie, calculateTimeDifference };