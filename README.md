# Liste de Tâches pour le Projet de Forum

## 1. Inscription et Connexion
- [ ] Créer un formulaire d'inscription avec les champs suivants :
  - Surnom
  - Âge
  - Genre
  - Prénom
  - Nom de famille
  - E-mail
  - Mot de passe
- [ ] Implémenter la logique de validation des données du formulaire d'inscription.
- [ ] Créer un système de connexion permettant l'authentification via pseudo ou e-mail et mot de passe.
- [ ] Implémenter la déconnexion accessible depuis n'importe quelle page.
- [ ] Gérer les sessions utilisateurs pour maintenir l'état de connexion.

## 2. Articles et Commentaires
- [ ] Créer une interface pour permettre aux utilisateurs de créer des publications.
- [ ] Implémenter la fonctionnalité de catégorisation des publications.
- [ ] Afficher les publications dans un flux sur la page d'accueil.
- [ ] Permettre aux utilisateurs de commenter les publications.
- [ ] Afficher les commentaires uniquement lorsque l'utilisateur clique sur une publication.

## 3. Messages Privés
- [ ] Créer une section de chat pour les messages privés.
- [ ] Afficher une liste d'utilisateurs en ligne/hors ligne, triée par dernier message ou par ordre alphabétique pour les nouveaux utilisateurs.
- [ ] Permettre l'envoi de messages privés aux utilisateurs en ligne.
- [ ] Implémenter une section pour afficher les messages passés d'une conversation.
- [ ] Charger les 10 derniers messages d'une conversation et permettre le chargement de 10 messages supplémentaires lors du défilement.
- [ ] Formater les messages pour inclure :
  - Date d'envoi
  - Nom d'utilisateur de l'expéditeur
- [ ] Assurer la réception en temps réel des messages via WebSockets.

## 4. Base de Données (SQLite)
- [ ] Concevoir le schéma de la base de données pour stocker les utilisateurs, publications, commentaires et messages privés.
- [ ] Implémenter les opérations CRUD (Créer, Lire, Mettre à jour, Supprimer) pour chaque entité.

## 5. Backend (Golang)
- [ ] Configurer le serveur Golang pour gérer les requêtes HTTP.
- [ ] Implémenter les WebSockets pour la communication en temps réel.
- [ ] Gérer les routes pour l'inscription, la connexion, la création de publications, les commentaires et les messages privés.

## 6. Frontend (HTML, CSS, JavaScript)
- [ ] Créer un fichier HTML unique pour l'application.
- [ ] Organiser les éléments de la page avec HTML.
- [ ] Styliser les éléments de la page avec CSS.
- [ ] Gérer les événements frontend avec JavaScript, y compris la navigation entre les différentes sections de l'application.
- [ ] Implémenter la logique pour gérer les WebSockets côté client.

## 7. Tests et Débogage
- [ ] Tester chaque fonctionnalité pour s'assurer qu'elle fonctionne comme prévu.
- [ ] Déboguer les problèmes rencontrés lors des tests.
- [ ] Effectuer des tests de performance pour s'assurer que l'application peut gérer plusieurs utilisateurs simultanément.

## 8. Documentation
- [ ] Documenter le code et les fonctionnalités de l'application.
- [ ] Rédiger un guide d'utilisation pour les utilisateurs finaux.
