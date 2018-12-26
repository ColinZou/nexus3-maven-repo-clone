# Nexus 3 Maven Repo Clone Guide
A simple guide for anyone want to copy a single maven repo from one nexus server to another. 
Or duplicate a maven repo on a same nexus server. I could put everything together into one piece latter.
## 0. Before You Start
Before you can continue, make sure you have access to you nexus server. This guide was tested in ubuntu 64bit server.
And you have python installed on the server.

Nexus 3 has its own storage structure. I basically was trying to parsing its .properties files for doing the job.  
## 1. Clone
### a). Create A Work Directory
    mkdir -p ~/nexus-clone
    cd ~/nexus-clone
Then put the tools clone-tool-linux-amd64 and generate-nexus-upload-command into this folder ~/nexus-clone.
    
    chmod +x ~/nexus-clone/clone-tool-linux-amd64
    chmod +x ~/nexus-clone/generate-nexus-upload-command
### b). List Artifacts Wants To Clone
Suppose your nexus data folder is /opt/nexus3-data, and you want to clone the maven repo "changeme" into ~/nexus-clone.
You must use a absolute path here.

    grep -l -r Bucket.repo-name=changeme /opt/nexus3-data/blobs/default > ~/nexus-clone/properties_list
    # check the total items we've chosen
    cat  ~/nexus-clone/properties_list | wc -l
### c). Use clone-tool-linux-amd64 To Copy jars/poms
    cat ~/nexus-clone/properties_list | sed 's|^|~/nexus-clone/clone-tool-linux-amd64 |g' > ~/nexus-clone/first_step
    sh ~/nexus-clone/first_step
You will see the jar/pom files are cloned. 

    find . -type f -name "*.jar"
    find . -type f -name "*.pom"
    
### d). Setup Your Nexus Repo
You need to setup your_maven_repository_id and your_maven_repo_url in generate-nexus-upload-command.

    post_fix = "-DrepositoryId=%s -Durl=%s"%("your_maven_repository_id","your_maven_repo_url")
    
You may find something about repositoryId from here if you have no idea about it.
https://www.baeldung.com/maven-deploy-nexus
### e). Upload 
You will see a file named call-generate-nexus-upload-command in your work folder. Run the following command.
    
    sh call-generate-nexus-upload-command | sh
    
You could check the output of the following command before you let shell run it. 
    
    sh call-generate-nexus-upload-command      
## Note
    # command for building the main.go
    env GOOS=linux GOARCH=amd64 go build -o clone-tool-linux-amd64 main.go 