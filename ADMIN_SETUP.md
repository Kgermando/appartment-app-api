# Configuration de l'Utilisateur Admin

## Fonctionnalité Admin Auto-Creation

Cette fonctionnalité a été ajoutée pour garantir qu'un utilisateur administrateur existe toujours dans le système.

### Comment ça fonctionne

1. **Vérification automatique** : À chaque tentative de connexion (`/api/auth/login`), le système vérifie automatiquement si un utilisateur avec le rôle "Administrator" existe.

2. **Création automatique** : Si aucun admin n'existe, le système crée automatiquement un utilisateur admin par défaut avec les credentials suivants :
   - **Email** : `admin@appartment-app.com`
   - **Téléphone** : `+243000000000`
   - **Mot de passe** : `Admin@123`
   - **Rôle** : `Administrator`
   - **Permission** : `ALL`
   - **Statut** : `Actif`

### Endpoints Disponibles

#### 1. Login avec vérification admin automatique
```http
POST /api/auth/login
Content-Type: application/json

{
  "identifier": "admin@appartment-app.com",
  "password": "Admin@123"
}
```

#### 2. Créer un admin manuellement
```http
POST /api/auth/create-admin
Content-Type: application/json

{
  "fullname": "Super Administrator",
  "email": "admin@example.com",
  "telephone": "+243123456789",
  "password": "SecurePassword123"
}
```

### Sécurité

⚠️ **IMPORTANT** : 
- Changez le mot de passe par défaut immédiatement après la première connexion
- Utilisez l'endpoint `/api/auth/change-password` pour modifier le mot de passe
- Assurez-vous que l'email et le téléphone par défaut ne sont pas accessibles par des tiers

### Test de la fonctionnalité

1. **Démarrer l'application**
2. **Faire une tentative de login** - cela déclenchera la création de l'admin si nécessaire
3. **Se connecter avec les credentials par défaut**
4. **Changer immédiatement le mot de passe**

### Logs

Quand un admin est créé automatiquement, vous verrez ces messages dans la console :
```
Utilisateur admin créé avec succès:
Email: admin@appartment-app.com
Téléphone: +243000000000
Mot de passe: Admin@123
⚠️  IMPORTANT: Changez le mot de passe par défaut après la première connexion!
```

### Modification du mot de passe

Une fois connecté, utilisez cet endpoint pour changer le mot de passe :

```http
PUT /api/auth/change-password?token=YOUR_JWT_TOKEN
Content-Type: application/json

{
  "old_password": "Admin@123",
  "password": "NewSecurePassword123",
  "password_confirm": "NewSecurePassword123"
}
```