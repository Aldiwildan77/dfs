# Discord File System (DFS)

**Description:**

Discord File System (DFS) is an innovative solution that transforms your Discord server into a collaborative file management system. It leverages Discord's robust API to create a virtual file-sharing network, allowing users to upload, organize, and retrieve files seamlessly across multiple channels or servers.  

With DFS, your community can enjoy features like structured folder management, file tagging, search capabilities, and cross-server file accessibility, making Discord a hub for efficient collaboration.  

**Features:**

- **Multi-Channel File Management:** Organize files into folders and subfolders within channels for easy navigation.  
- **Cross-Server File Access:** Share and access files across interconnected servers.  
- **Advanced Search:** Locate files quickly with smart filters and tags.  
- **Secure Storage:** Ensure privacy and data integrity with robust encryption and audit logs.  

DFS is perfect for teams, communities, and creators looking to optimize their Discord experience for productivity and collaboration. 🚀

## Usage

### Save a file to the DFS

```curl
curl --location 'http://localhost:8080/api/file/save' \
--header 'Content-Type: application/json' \
--data '{
    "guild_id": "919965514027547782",
    "channel_id": "919965514027547782",
    "message_id": "919965514027547782",
    "url": "https://cdn.discordapp.com/attachments/919965514027547782/1319544033684754475/2769762.png?ex=676852e5&is=67670165&hm=a50f68cbf89ce5fa26b84d238fec33362afdaac8c097ac899dd60e12cb4670eb&"
}'
```

### Retrieve a file from the DFS

```curl
curl --location 'http://localhost:8080/api/file/get' \
--header 'Content-Type: application/json' \
--data '{
    "url": "https://cdn.discordapp.com/attachments/919965514027547782/1319544033684754475/2769762.png?ex=676852e5&is=67670165&hm=a50f68cbf89ce5fa26b84d238fec33362afdaac8c097ac899dd60e12cb4670eb&"
}'
```

## Checklist

- [x] Save uploaded file to the database
- [x] Access uploaded file from the database
- [x] Rotating expired file
- [ ] Direct upload to Discord
- [ ] Manage folder and file
- [ ] Search file by name, tag, and date

## Authors

- [Muhammad Wildan Aldiansyah](https://github.com/Aldiwildan77)
