// Subdomain Enumeration and Filtering

// Note
// This is just for demonstration 
// This script will probably failed because I have not released
// Some of my own tools publicly yet and these filepaths are mine


//Important File Paths
@amasspassive = "/opt/tools/Amass/passive.ini"
@resolvers = "/opt/wordlists/resolvers.txt"
@dnsmicro = "/opt/wordlists/dnsbrute.txt"



//Create Temporary Inscope File
echo @rootsubs{file} #as:@inscopefile

//Create Temporary OutofScope File
echo @outscope{!file}  #as:@outscopefile

// ####---- Passive Subdomain Enumeration ---

// assetfinder subdomains
assetfinder -subs-only #from:@rootsubs #as:@passivesubs{unique}

//Passive Amass Scan 
amass enum -df @inscopefile -blf @outscopefile -config @amasspassive -o @outfile #as:@passivesubs{unique}

// Subfinder Enumeration
subfinder -all -dL @inscopefile -nc -silent #as:@passivesubs{unique}

// Findomain Subdomains
findomain -f @inscopefile -o @outfile #as:@passivesubs{unique}

// Filter inscope and outscope  domains from passive enum
subops filter --df @rootsubs{file} --ef @outscope{!file}  #from:@passivesubs #as:@passive{unique}




// ####--- Active Subdomain Enumeration

// DNS BruteForce using PureDNS
puredns bruteforce @dnsmicro -r @resolvers -w @outfile @z #for:@rootsubs:@z #as:@activesubs{unique}

// Get Subdomains from TLS Certificates
cero -p 443,4443,8443,10443 -c 1000  #from:@passive #as:@activesubs{unique}

// Get Subdomains From CNAMEs
dnsx -retry 3 -cname -l @passivesubs{file} -o @outfile #as:@activesubs{unique}

// Filter inscope and outscope  domains from active enum
subops filter --df @rootsubs{file} --ef @outscope{!file} #from:@activesubs #as:@active{unique}


// Merge Results From Active and Passive Sources


//Check If Active Sub Enum Was Useful
bninja diff @active{file} @passive{file} -1 #as:@uniqueactivesubs #notifylen{Total Subdomain Found From Active Enum but Not Passive :}

//Merge passive and activesubs
cat @passive{file} @active{file} | bninja uniq #as:@allsubs

//Check For Subdomain Takeover If in Scope
subjack -w @allsubs{file} -t 100 -timeout 30 -o @outfile -ssl


// ###---- Filtering Found Subdomains

//Resolve All Passive Found Subdomains
rusolver -r @resolvers -i #from:@allsubs #as:@resolved

//Filter Ips Pointing to Private Ip Addresses
rusolverfilter  -sub-only #from:@resolved #as:@filtered  #notifylen{Total Resolved Subdomains: }

//Subdomins Hosting Web Services API,Web Page etc
httpx-pd -silent -t 100 -retries 2 #from:@filtered  #as:@webraw #notifylen{Total Subdomains With Web Services :}

//remove https from httpx output
cut -d "/" -f 3 #from:@webraw #as:@webservices

//Subdomains Not Hosting Web Services 
bninja diff @filtered{file} @webservices{file} -1 #as:@nonwebsubs #notifylen{Total Non Web Subs }

//Port Scan All Non Web
naabu -ec -ep 21,22,80,443 -verify -nmap-cli "nmap -sC -sV" -o @outfile -c 50 -list @nonwebsubs{file} #as:@nonwebportscan #notify:{Non Web Port Scan Results}

