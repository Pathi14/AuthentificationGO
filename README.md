# AuthentificationGO

## Organisation des dossiers

[Guide sur l'organisation des dossiers en Go](https://medium.com/@smart_byte_labs/organize-like-a-pro-a-simple-guide-to-go-project-folder-structures-e85e9c1769c2)

### Architecture Domain-Driven Design (DDD)

L'application utilise une architecture **Domain-Driven Design (DDD)**. Cette approche consiste à diviser l'application en domaines ou contextes délimités (bounded contexts), où chaque domaine gère ses propres couches (modèles, répositories, services).

Avantages de l'approche DDD :
- Isolation de la logique métier
- Organisation claire et modulaire
- Favorise la maintenabilité et l'évolutivité

## Structure du projet

```
├── cmd/                  # Point d'entrée de l'application
├── internal/             # Code propre à l'application
│   ├── user/            # Domaine utilisateur (models, services, repositories)
|   ├── infrastructure/  # Infrastructure de l'application (bases de données, middlewares)
|   ├── middleware/      # Middlewares
│   └── ...              # Autres domaines
├── tests/                # Tests d'intégration
└── main.go               # Fichier principal
```

## Prérequis

- Go (>= 1.20)
- PostgreSQL

## Installation

1. Clonez le dépôt :

```bash
git clone https://github.com/Pathi14/AuthentificationGO.git
cd AuthentificationGO
```

2. Créez un fichier `.env` à la racine :

```env
DB_HOST=localhost
DB_PORT=5433
DB_USER=admin
DB_PASSWORD=secret
DB_NAME=authentificationgo
JWT_SECRET=Up7j~wFP{y2?cqvk}x'W)X
```

3. Installez les dépendances :

```bash
go mod tidy
```

## Lancer l'application

1. Démarrer la base de données (PostgreSQL) :

Assurez-vous que votre base de données est accessible et configurée correctement.
Pour démarrer la base de données PostgreSQL en développement ainsi que la base de données utilisée pour les tests, il vous suffit de lancer Docker Compose :
```bash
docker compose up -d
```

2. Exécutez l'application :

```bash
go run main.go
ou
go run .
```

## Utilisation de l'API

### Créer un utilisateur

```bash
POST /44df37e7-fe2a-404f-917b-399f5c5ffd12/register
{
  "name": "John Doe",
  "age": 30,
  "mobile_number": "1234567890",
  "email": "john@example.com",
  "password": "Secret123"
}

```

### Se connecter

```bash
POST /44df37e7-fe2a-404f-917b-399f5c5ffd12/login
{
  "email": "john@example.com",
  "password": "Secret123"
}
```

### Mot de passe oublié

```bash
POST /44df37e7-fe2a-404f-917b-399f5c5ffd12/forgot-password
{
  "email": "john@example.com"
}
```

### Réinitialiser le mot de passe

```bash
POST /44df37e7-fe2a-404f-917b-399f5c5ffd12/reset-password
{
  "token": "reset-password-token",
  "new_password": "NewSecret123"
}
```

### Voir le profil de l'utilisateur (Route protégée)

```bash
GET /44df37e7-fe2a-404f-917b-399f5c5ffd12/me
```

### Se déconnecter (Route protégée)

```bash
POST /44df37e7-fe2a-404f-917b-399f5c5ffd12/logout
```

### Raffraichir le jeton d'accès (Route protégée)

```bash
POST /44df37e7-fe2a-404f-917b-399f5c5ffd12/refresh
Authorization: Bearer your-refresh-token
```

### 

## Exécuter les tests

Pour exécuter l'ensemble des tests :

```bash
go test ./tests
```

## Ressources utiles

- [Tutoriel : backend en Go avec PostgreSQL](https://medium.com/@minduladilthushan/building-a-simple-backend-server-in-go-with-postgresql-and-testing-with-postman-92f7796f696c)

- [Guide sur l'organisation des dossiers en Go](https://medium.com/@smart_byte_labs/organize-like-a-pro-a-simple-guide-to-go-project-folder-structures-e85e9c1769c2)