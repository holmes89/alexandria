(defproject alexandria "0.1.0-SNAPSHOT"
  :dependencies [[org.clojure/clojure "1.10.1"]
                 [org.clojure/clojurescript "1.10.520"
                  :exclusions [com.google.javascript/closure-compiler-unshaded
                               org.clojure/google-closure-library]]
                 [thheller/shadow-cljs "2.8.65"]
                 [cljs-ajax "0.8.0"]
                 [day8.re-frame/http-fx "0.1.6"]
                 [reagent "0.8.1"]
                 [re-frame "0.10.9"]
                 [secretary "1.2.3"]]

  :plugins [
            ]

  :min-lein-version "2.5.3"

  :jvm-opts ["-Xmx1G"]

  :source-paths ["src/clj" "src/cljs"]

  :clean-targets ^{:protect false} ["resources/public/js/compiled" "target"]


  :aliases {"dev"        ["with-profile" "dev" "run" "-m" "shadow.cljs.devtools.cli" "watch" "app"]
            "prod"       ["with-profile" "prod" "run" "-m" "shadow.cljs.devtools.cli" "release" "app"]
            "karma-once" ["with-profile" "prod" "do"
                          ["clean"]
                          ["run" "-m" "shadow.cljs.devtools.cli" "compile" "karma-test"]
                          ["shell" "karma" "start" "--single-run" "--reporters" "junit,dots"]]}


  :profiles
  {:dev
   {:dependencies [[binaryage/devtools "0.9.10"]
                   [re-frisk "0.5.4.1"]]}

   :prod { }
   })
