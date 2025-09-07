# Dashboard API - Guide d'utilisation

## Vue d'ensemble

Le dashboard fournit des statistiques complètes et des analyses pour la gestion des appartements et des finances. Il inclut des filtres avancés, notamment la possibilité de filtrer par manager ou de voir les données pour tous les managers.

## Endpoints disponibles

### 1. Dashboard Principal - Statistiques générales
```
GET /api/dashboard/stats
```

**Paramètres de requête (optionnels) :**
- `manager_uuid` : Filtrer par un manager spécifique
- `start_date` : Date de début (format: YYYY-MM-DD)
- `end_date` : Date de fin (format: YYYY-MM-DD)

**Exemples d'utilisation :**
```bash
# Dashboard global (tous les managers)
GET /api/dashboard/stats

# Dashboard pour un manager spécifique
GET /api/dashboard/stats?manager_uuid=123e4567-e89b-12d3-a456-426614174000

# Dashboard avec filtres de date
GET /api/dashboard/stats?start_date=2024-01-01&end_date=2024-12-31

# Dashboard pour un manager avec filtres de date
GET /api/dashboard/stats?manager_uuid=123e4567-e89b-12d3-a456-426614174000&start_date=2024-01-01&end_date=2024-12-31
```

**Réponse :**
```json
{
  "status": "success",
  "message": "Dashboard stats retrieved successfully",
  "data": {
    "total_apartments": 150,
    "available_apartments": 25,
    "occupied_apartments": 120,
    "maintenance_apartments": 5,
    "total_income_cdf": 500000.0,
    "total_income_usd": 2500.0,
    "total_expense_cdf": 150000.0,
    "total_expense_usd": 750.0,
    "net_balance_cdf": 350000.0,
    "net_balance_usd": 1750.0,
    "monthly_revenue_target": 75000.0,
    "actual_monthly_revenue": 68000.0,
    "revenue_percentage": 90.67,
    "top_apartments_by_revenue": [...],
    "manager_stats": [...]
  }
}
```

### 2. Tendances mensuelles
```
GET /api/dashboard/trends
```

**Paramètres de requête (optionnels) :**
- `manager_uuid` : Filtrer par un manager spécifique
- `months` : Nombre de mois à afficher (défaut: 12)

**Exemples d'utilisation :**
```bash
# Tendances des 12 derniers mois pour tous les managers
GET /api/dashboard/trends

# Tendances des 6 derniers mois pour un manager spécifique
GET /api/dashboard/trends?manager_uuid=123e4567-e89b-12d3-a456-426614174000&months=6
```

### 3. Comparaison entre managers
```
GET /api/dashboard/managers
```

**Paramètres de requête (optionnels) :**
- `start_date` : Date de début (format: YYYY-MM-DD)
- `end_date` : Date de fin (format: YYYY-MM-DD)

**Exemples d'utilisation :**
```bash
# Comparaison de tous les managers
GET /api/dashboard/managers

# Comparaison avec filtres de date
GET /api/dashboard/managers?start_date=2024-01-01&end_date=2024-12-31
```

**Réponse :**
```json
{
  "status": "success",
  "message": "Manager comparison retrieved successfully",
  "data": [
    {
      "manager_uuid": "123e4567-e89b-12d3-a456-426614174000",
      "manager_name": "Jean Dupont",
      "total_apartments": 25,
      "available_apartments": 5,
      "occupied_apartments": 20,
      "total_income_cdf": 200000.0,
      "total_income_usd": 1000.0,
      "total_expense_cdf": 50000.0,
      "total_expense_usd": 250.0,
      "net_balance_cdf": 150000.0,
      "net_balance_usd": 750.0,
      "monthly_revenue_target": 30000.0
    }
  ]
}
```

### 4. Performance des appartements
```
GET /api/dashboard/apartments/performance
```

**Paramètres de requête (optionnels) :**
- `manager_uuid` : Filtrer par un manager spécifique
- `start_date` : Date de début (format: YYYY-MM-DD)
- `end_date` : Date de fin (format: YYYY-MM-DD)
- `limit` : Nombre maximum de résultats (défaut: 20)

**Exemples d'utilisation :**
```bash
# Top 20 appartements par revenus (tous managers)
GET /api/dashboard/apartments/performance

# Top 10 appartements pour un manager spécifique
GET /api/dashboard/apartments/performance?manager_uuid=123e4567-e89b-12d3-a456-426614174000&limit=10
```

## Données retournées

### Statistiques principales
- **total_apartments** : Nombre total d'appartements
- **available_apartments** : Appartements disponibles
- **occupied_apartments** : Appartements occupés
- **maintenance_apartments** : Appartements en maintenance
- **total_income_cdf/usd** : Revenus totaux en CDF et USD
- **total_expense_cdf/usd** : Dépenses totales en CDF et USD
- **net_balance_cdf/usd** : Balance nette en CDF et USD
- **monthly_revenue_target** : Objectif de revenus mensuels
- **actual_monthly_revenue** : Revenus mensuels actuels
- **revenue_percentage** : Pourcentage de réalisation des objectifs

### Top appartements par revenus
- Classement des appartements les plus rentables
- Inclut nom, numéro, loyer mensuel, revenus totaux, statut
- Nom du manager responsable

### Statistiques par manager
- Statistiques détaillées pour chaque manager
- Nombre d'appartements gérés par statut
- Revenus, dépenses et balance nette
- Objectifs de revenus mensuels

## Cas d'utilisation

### 1. Dashboard d'administration (voir tous les managers)
```javascript
// Récupérer les statistiques globales
fetch('/api/dashboard/stats')
  .then(response => response.json())
  .then(data => {
    console.log('Statistiques globales:', data.data);
    console.log('Stats par manager:', data.data.manager_stats);
  });

// Comparaison entre managers
fetch('/api/dashboard/managers')
  .then(response => response.json())
  .then(data => {
    console.log('Comparaison managers:', data.data);
  });
```

### 2. Dashboard manager spécifique
```javascript
const managerUUID = '123e4567-e89b-12d3-a456-426614174000';

// Statistiques pour un manager spécifique
fetch(`/api/dashboard/stats?manager_uuid=${managerUUID}`)
  .then(response => response.json())
  .then(data => {
    console.log('Stats du manager:', data.data);
  });

// Performance des appartements du manager
fetch(`/api/dashboard/apartments/performance?manager_uuid=${managerUUID}`)
  .then(response => response.json())
  .then(data => {
    console.log('Performance appartements:', data.data);
  });
```

### 3. Analyses avec filtres de date
```javascript
const startDate = '2024-01-01';
const endDate = '2024-12-31';

// Statistiques pour une période spécifique
fetch(`/api/dashboard/stats?start_date=${startDate}&end_date=${endDate}`)
  .then(response => response.json())
  .then(data => {
    console.log('Stats période:', data.data);
  });

// Tendances mensuelles
fetch(`/api/dashboard/trends?months=12`)
  .then(response => response.json())
  .then(data => {
    console.log('Tendances:', data.data);
  });
```

## Filtres disponibles

1. **Par Manager** : `manager_uuid`
   - Filtre toutes les données par un manager spécifique
   - Si omis, affiche les données de tous les managers

2. **Par Date** : `start_date` et `end_date`
   - Filtre les transactions (caisses) par période
   - Format : YYYY-MM-DD
   - Affecte les calculs de revenus/dépenses

3. **Limite de résultats** : `limit`
   - Limite le nombre de résultats retournés
   - Utile pour les listes comme "top appartements"

## Sécurité et permissions

- Authentification requise pour tous les endpoints
- Les managers ne peuvent voir que leurs propres données
- Les administrateurs/superviseurs peuvent voir toutes les données
- Validation des paramètres d'entrée
- Protection contre l'injection SQL

## Performance

- Requêtes optimisées avec jointures
- Index sur les colonnes de filtrage
- Limitation des résultats pour éviter les surcharges
- Cache recommandé pour les données fréquemment consultées
