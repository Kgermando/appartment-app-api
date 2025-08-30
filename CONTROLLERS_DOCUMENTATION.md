# Documentation des Contrôleurs et Système de Conversion

## 🏠 Contrôleur Appartment

Le contrôleur `appartment.controller.go` gère toutes les opérations CRUD pour les appartements.

### Endpoints disponibles :
- `GET /appartments/paginated/:manager_uuid` - Pagination avec filtrage par manager
- `GET /appartments/paginated` - Pagination générale
- `GET /appartments` - Tous les appartements
- `GET /appartments/manager/:manager_uuid` - Appartements par manager
- `GET /appartments/:uuid` - Un appartement spécifique
- `POST /appartments` - Créer un appartement
- `PUT /appartments/:uuid` - Mettre à jour un appartement
- `DELETE /appartments/:uuid` - Supprimer un appartement

## 💰 Contrôleur Caisse

Le contrôleur `caisse.controller.go` gère les entrées et sorties financières avec système de conversion.

### Endpoints disponibles :

#### Opérations CRUD de base :
- `GET /caisses/paginated/:appartment_uuid` - Pagination par appartement
- `GET /caisses/paginated` - Pagination générale
- `GET /caisses` - Toutes les caisses
- `GET /caisses/appartment/:appartment_uuid` - Caisses par appartement
- `GET /caisses/:uuid` - Une caisse spécifique
- `POST /caisses` - Créer une entrée de caisse
- `PUT /caisses/:uuid` - Mettre à jour une caisse
- `DELETE /caisses/:uuid` - Supprimer une caisse

#### Fonctionnalités financières avancées :

**1. Balance par appartement :**
```
GET /caisses/balance/:appartment_uuid
```
Retourne :
- Total Income CDF/USD
- Total Expense CDF/USD
- Balance nette CDF/USD

**2. Totaux globaux :**
```
GET /caisses/totals/global
```
Retourne :
- Totaux Income/Expense en CDF et USD
- Conversions croisées (CDF en USD, USD en CDF)
- Grand totaux combinés
- Balances nettes
- Taux de change actuels

**3. Totaux par manager :**
```
GET /caisses/totals/manager/:manager_uuid
```
Retourne les mêmes informations que les totaux globaux mais filtrés par manager.

**4. Conversion de devises :**
```
POST /caisses/convert
```
Corps de la requête :
```json
{
  "amount": 100.50,
  "from_currency": "USD",
  "to_currency": "CDF"
}
```

## 🔄 Système de Conversion de Devises

### Utilitaire de conversion (`utils/currency.go`)

**Fonctionnalités :**
- Conversion automatique USD ↔ CDF
- Récupération des taux depuis API externe
- Taux de change par défaut en cas d'échec de l'API
- Support pour d'autres devises (extensible)

**Fonctions principales :**
- `ConvertCurrency(amount, from, to)` - Conversion générale
- `ConvertUSDToCDF(amount)` - USD vers CDF
- `ConvertCDFToUSD(amount)` - CDF vers USD
- `GetCurrentExchangeRate(from, to)` - Taux actuel

**Configuration par défaut :**
- USD vers CDF : 2700 (modifiable)
- CDF vers USD : 0.00037 (modifiable)

### API de taux de change

Le système utilise l'API gratuite `exchangerate-api.com` par défaut, mais peut être configuré pour utiliser d'autres APIs.

## 📊 Contrôleur Exchange (Bonus)

Le contrôleur `exchange.controller.go` gère les taux de change en base de données.

### Endpoints :
- `GET /exchange/active` - Taux actifs
- `GET /exchange/paginated` - Pagination des taux
- `GET /exchange/:uuid` - Un taux spécifique
- `POST /exchange` - Créer un taux
- `PUT /exchange/:uuid` - Mettre à jour un taux
- `DELETE /exchange/:uuid` - Supprimer un taux
- `POST /exchange/sync` - Synchroniser depuis l'API

### Modèle ExchangeRate

```go
type ExchangeRate struct {
    UUID         string     `json:"uuid"`
    FromCurrency string     `json:"from_currency"`
    ToCurrency   string     `json:"to_currency"`
    Rate         float64    `json:"rate"`
    Source       string     `json:"source"` // manual, api, automatic
    IsActive     bool       `json:"is_active"`
    ValidFrom    time.Time  `json:"valid_from"`
    ValidTo      *time.Time `json:"valid_to,omitempty"`
    UpdatedByUUID string    `json:"updated_by_uuid"`
    // Relations et timestamps...
}
```

## 🛠 Utilisation

### Exemple de calcul de totaux globaux :

La réponse de `GET /caisses/totals/global` inclut :

```json
{
  "status": "success",
  "message": "Global totals retrieved successfully",
  "data": {
    "income_totals": {
      "cdf_total": 15000000,
      "usd_total": 2500,
      "cdf_in_usd": 5555.56,
      "usd_in_cdf": 6750000,
      "grand_total_cdf": 21750000,
      "grand_total_usd": 8055.56
    },
    "expense_totals": {
      "cdf_total": 8000000,
      "usd_total": 1200,
      "cdf_in_usd": 2962.96,
      "usd_in_cdf": 3240000,
      "grand_total_cdf": 11240000,
      "grand_total_usd": 4162.96
    },
    "net_balances": {
      "net_balance_cdf": 7000000,
      "net_balance_usd": 1300,
      "grand_net_balance_cdf": 10510000,
      "grand_net_balance_usd": 3892.60
    },
    "exchange_rates": {
      "usd_to_cdf": 2700,
      "cdf_to_usd": 0.00037
    }
  }
}
```

### Validation des types de caisse :

Les caisses acceptent uniquement deux types :
- `"Income"` - Pour les entrées (loyers, revenus)
- `"Expense"` - Pour les sorties (maintenance, frais)

### Devises supportées :

Actuellement supportées :
- `"USD"` - Dollar américain
- `"CDF"` - Franc congolais

Le système est extensible pour ajouter d'autres devises.

## 📈 Avantages du système

1. **Calculs automatiques** : Conversion automatique entre devises
2. **Historique** : Stockage des taux de change historiques
3. **Flexibilité** : API externe + taux par défaut
4. **Performance** : Requêtes optimisées avec COALESCE
5. **Validation** : Vérification stricte des types et devises
6. **Pagination** : Support de la pagination pour toutes les listes
7. **Recherche** : Fonctionnalités de recherche intégrées
