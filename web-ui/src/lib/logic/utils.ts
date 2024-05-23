import QRCode from 'qrcode'
import { parseISO, format } from 'date-fns';

function checkMail(mail: string): boolean {
  return /^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,6}$/.test(mail);
}

function setCookie(data: string) {
  document.cookie = data;
}

function betterTime(lastUse: string) {
  return parseISO(lastUse).toLocaleString();
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

function formatUserAgent(userAgent: string): string {
  // Inizia costruendo una stringa vuota per la descrizione
  let description = "";

  // Analizza l'User-Agent per identificare il sistema operativo
  const osMatch = userAgent.toLowerCase().includes("windows");
  if (osMatch) {
      description += "Windows ";
  } else if (userAgent.toLowerCase().includes("macintosh") || userAgent.toLowerCase().includes("mac os x")) {
      description += "Mac OS X ";
  } else if (userAgent.toLowerCase().includes("iphone") || userAgent.toLowerCase().includes("ipad") || userAgent.toLowerCase().includes("ipod")) {
      description += "iOS ";
  } else if (userAgent.toLowerCase().includes("android")) {
      description += "Android ";
  } else if (userAgent.toLowerCase().includes("linux")) {
      description += "Linux ";
  } else {
      description += "Unknow OS ";
  }

  description += " (";

  // Analizza l'User-Agent per identificare il browser
  const browserMatch = userAgent.toLowerCase().includes("chrome");
  if (browserMatch) {
      description += "Chrome";
  } else if (userAgent.toLowerCase().includes("firefox")) {
      description += "Firefox";
  } else if (userAgent.toLowerCase().includes("safari")) {
      description += "Safari";
  } else if (userAgent.toLowerCase().includes("opera")) {
      description += "Opera";
  } else if (userAgent.toLowerCase().includes("msie") || userAgent.toLowerCase().includes("trident")) {
      description += "Internet Explorer";
  } else {
      description += "Unknown Browser";
  }

  return description + ")";
}

export { checkMail, formatDate, setCookie, generateQRCode, betterTime, getCookie, formatUserAgent, calculateTimeDifference };