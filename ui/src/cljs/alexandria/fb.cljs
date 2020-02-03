(ns alexandria.fb
  (:require ["firebase/app" :as firebase]
            ["firebase/auth"]))

(goog-define apikey "")
(goog-define projectid "")
(goog-define authdomain "")


;;(def login-button

(defn sign-in-with-github
  []
  (let [provider (firebase/auth.GithubAuthProvider.)]
    (.signInWithPopup (firebase/auth) provider)))


(defn on-auth-state-changed
  []
  (.onAuthStateChanged
    (firebase/auth)
    (fn
      [user]
      (if user
        (let [uid (.-uid user)
              display-name (.-displayName user)
              photo-url (.-photoURL user)
              email (.-email user)]
          (do
            ;; TODO dispatch
            (reset! state/user {:photo-url photo-url
                                :display-name display-name
                                :email email})))
        (reset! state/user nil)))))

(defn sign-out
  []
  (.signOut (firebase/auth)))

(defn firebase-init
  []
  (firebase/initializeApp
    #js {:apiKey "your-api-key"
         :authDomain "your-auth-domain"
         :projectId "your-project-id"})
  (on-auth-state-changed))
