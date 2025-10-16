# LinkedIn Post - HiveSSH

## 🚀 Presentando HiveSSH: Tu Controlador de SSH Múltiple

¿Cansado de gestionar múltiples conexiones SSH y ejecutar comandos remotos uno por uno? 

Te presento **HiveSSH** 🐝 - una herramienta CLI desarrollada en Go que simplifica la administración de servidores remotos.

### ✨ Características principales:

🔹 **Gestión centralizada** - Almacena y organiza todos tus servidores SSH en un solo lugar
🔹 **Ejecución en grupo** - Ejecuta comandos en múltiples servidores simultáneamente
🔹 **Agrupación inteligente** - Organiza tus servidores por clusters, entornos o proyectos
🔹 **Monitoreo de estado** - Verifica rápidamente qué servidores están disponibles
🔹 **Múltiples formatos de salida** - JSON, CSV y más

### 💻 Casos de uso:

✅ Gestionar clusters de Kubernetes
✅ Despliegues en múltiples entornos
✅ Administración de infraestructura distribuida
✅ Mantenimiento de servidores en lote

### 🛠️ Ejemplo de uso:

```bash
# Añadir un servidor
hivessh join master 192.168.0.23 --user root --description "Master K8s"

# Crear un grupo
hivessh group create production

# Ejecutar comando en grupo
hivessh run "kubectl get pods" --group production

# Listar todos los servidores
hivessh list --output json
```

### 🔧 Tecnologías:

- Go 1.24
- Cobra CLI framework
- SSH/SFTP libraries
- JSON-based storage

Este proyecto nace de la necesidad de simplificar operaciones diarias en infraestructuras con múltiples servidores. ¡Perfecto para DevOps, SysAdmins y desarrolladores!

🔗 Repositorio: https://github.com/Izangildev/hiveSSH

¿Qué opinas? ¿Te gustaría probar HiveSSH en tu infraestructura? 

#DevOps #SSH #Golang #CLI #OpenSource #SysAdmin #Infrastructure #Automation #Kubernetes #CloudComputing

---

## English Version

## 🚀 Introducing HiveSSH: Your Multiple SSH Controller

Tired of managing multiple SSH connections and executing remote commands one by one?

Meet **HiveSSH** 🐝 - a CLI tool built in Go that simplifies remote server management.

### ✨ Key Features:

🔹 **Centralized Management** - Store and organize all your SSH servers in one place
🔹 **Group Execution** - Run commands on multiple servers simultaneously
🔹 **Smart Grouping** - Organize your servers by clusters, environments, or projects
🔹 **Status Monitoring** - Quickly check which servers are available
🔹 **Multiple Output Formats** - JSON, CSV, and more

### 💻 Use Cases:

✅ Manage Kubernetes clusters
✅ Multi-environment deployments
✅ Distributed infrastructure administration
✅ Batch server maintenance

### 🛠️ Usage Example:

```bash
# Add a server
hivessh join master 192.168.0.23 --user root --description "Master K8s"

# Create a group
hivessh group create production

# Run command on group
hivessh run "kubectl get pods" --group production

# List all servers
hivessh list --output json
```

### 🔧 Technologies:

- Go 1.24
- Cobra CLI framework
- SSH/SFTP libraries
- JSON-based storage

This project was born from the need to simplify daily operations on infrastructures with multiple servers. Perfect for DevOps, SysAdmins, and developers!

🔗 Repository: https://github.com/Izangildev/hiveSSH

What do you think? Would you like to try HiveSSH in your infrastructure?

#DevOps #SSH #Golang #CLI #OpenSource #SysAdmin #Infrastructure #Automation #Kubernetes #CloudComputing
