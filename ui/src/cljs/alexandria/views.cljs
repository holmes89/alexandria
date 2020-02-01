(ns alexandria.views
  (:require
   [re-frame.core :as re-frame]
   [alexandria.subs :as subs]
   [alexandria.events :as events]
   [react-pdf :as pdf]
   [clojure.string :as str]))


;; shared components
(defn button [text on-click]
  [:a.button.is-light
   {:on-click on-click}
   text])


(defn navbar []
  [:nav.navbar.is-light {:role "navigation" :aria-label "main navigation"}
   [:div.navbar-brand
    [:div.navbar-item
     [:span
      "Alexandria"
      [:img {:src "/assets/alexandria.png"}]]]]])

;; home
(defn file-upload-name []
  (peek (str/split (.-value (.getElementById js/document "file-upload")) "\\")))
(defn upload-form [] (js/FormData. (.getElementById js/document "upload")))
(defn set-upload-name [name] (re-frame/dispatch [::events/update-upload-file-name name]))
(defn submit-upload []
  (re-frame/dispatch [::events/upload-book (upload-form)])
  (re-frame/dispatch [::events/hide-upload-modal])
  (set-upload-name ""))

(defn upload-modal []
  (let [show (re-frame/subscribe [::subs/is-upload-showing?])
        upload-name (re-frame/subscribe [::subs/upload-file-name])]
    (if @show
      [:div.modal.is-active
       [:div.modal-background]
       [:div.modal-card
        [:header.modal-card-head
         [:p.modal-card-title "Upload"]
         [:button.delete {:on-click #(re-frame/dispatch [::events/hide-upload-modal])}]]
        [:section.modal-card-body
         [:form#upload
          [:div.field
           [:label.label "Name"]
           [:div.control
            [:input.input {:type "text" :name "name"}]]]
          [:div.field
           [:label.file-label
            [:input#file-upload.file-input {:type "file" :name "file" :on-change #(set-upload-name (file-upload-name))}]
            [:span.file-cta
             [:span.file-icon
              [:i.fas.fa-upload]]
             [:span.file-label "Choose a file..."]]
            [:span.file-name @upload-name]]]]]
        [:footer.modal-card-foot
         [:button.button.is-success {:on-click submit-upload}"Submit"]
         [:button.button {:on-click #(re-frame/dispatch [::events/hide-upload-modal])} "Cancel"]]]]
      [:div])))

(defn doc-icon
  [type]
  (if (= type "book")
    [:i.fas.fa-book]
    [:i.fas.fa-file-alt]))



(defn doc-card
  [{:keys [id display_name type description]}]
  [:div.card
   [:header.card-header
    [:p.card-header-icon
     [:span.icon [doc-icon type]]]
    [:p.card-header-title display_name]]
   [:div.card-content {:style {:text-align "center"}}
    [:img {:src (str "http://read.jholmestech.com/assets/covers/" id ".jpg") :style {:max-width "300px"}}]]
   [:footer.card-footer
    [:a.card-footer-item {:href (str "#/documents/" id)} [:i.fas.fa-book-open] "Read"]]])

(defn doc-card-grid []
  (let [docs (re-frame/subscribe [::subs/docs])]
    (fn []
      [:div.columns.is-mobile.is-multiline
       (for [doc @docs]
         ^{:key (:id doc)}
         [:div.column.is-one-quarter
          [doc-card doc]])])))


(defn add-button []
  [:div.is-pulled-right
   [:button.button.is-rounded.is-info {:on-click #(re-frame/dispatch [::events/show-upload-modal])}
    [:span.icon
     [:i.fas.fa-plus]]
    [:span "Add"]]])

(defn authenticated-body []
  (re-frame/dispatch [::events/get-documents])
  [:div
   [:section.section
    [:div.container
     [:div.columns.is-mobile
      [:div.column.is-6.is-offset-3
       [:input.input {:type "text" :placeholder "Search"}]]
      [:div.column.is-1.is-offset-2 (add-button)]]]]
   [:section.section
    [:div.container
     [doc-card-grid]
     (upload-modal)]]])

(defn unauthenticated-body []
  [:div
   [:h1.main-title "Alexandria" ]
   [:div "A self managed library of documents"]])

(defn home-panel []
  (let [name (re-frame/subscribe [::subs/name])]
    [:div
     (navbar)
     (authenticated-body)]))


;; read
(defn pdf-page [num]
  (let [page-num (re-frame/subscribe [::subs/doc-page])
        zoom (re-frame/subscribe [::subs/doc-zoom])]
    [:> pdf/Page {:pageNumber @page-num :scale @zoom :renderAnnotationLayer false}]))

(defn zoom-in []
  [:a {:on-click #(re-frame/dispatch [::events/zoom-in])}
   [:i.fas.fa-search-plus]])

(defn zoom-out []
  [:a { :on-click #(re-frame/dispatch [::events/zoom-out])}
   [:i.fas.fa-search-minus]])

(defn read-bar
  [{:keys [display_name path name]}]
  [:nav#read-bar.navbar.is-dark
   [:div.navbar-menu
    [:div.navbar-start
     [:div.navbar-item
      [:a {:href "/#/"} [:i.fas.fa-arrow-left]]]]
    [:div.navbar-start.centered
     [:div.navbar-item
      [:h3.book-title display_name]]]
    [:div.navbar-end
     [:a {:href path :download name :target "_blank"} [:i.fas.fa-download]]
     (zoom-in)
     (zoom-out)]]])

(defn next-page []
  [:a.page-turn {:on-click #(re-frame/dispatch [::events/next-page])}
   [:i.fas.fa-arrow-right]])

(defn prev-page []
  [:a.page-turn {:on-click #(re-frame/dispatch [::events/prev-page])}
   [:i.fas.fa-arrow-left]])

(defn pdf-reader [src]
  [:> pdf/Document {:file src}
   (pdf-page 1)])

(defn read-section []
  (let [doc (re-frame/subscribe [::subs/active-doc])]
    [:div
     (read-bar @doc)
     [:div
      [:div.columns.is-gapless
       [:div.column.is-1
        (prev-page)]
       [:div#doc.column.is-10
        (let [src (:path @doc)]
          (if src
            [pdf-reader (:path @doc)]))]
       [:div.column.is-1
        (next-page)]]]]))


(defn read-panel []
  [:div
   (navbar)
   (read-section)])



;; main

(defn- panels [panel-name]
  (case panel-name
    :home-panel [home-panel]
    :read-panel [read-panel]
    [:div]))

(defn show-panel [panel-name]
  [panels panel-name])

(defn main-panel []
  (let [active-panel (re-frame/subscribe [::subs/active-panel])]
    [show-panel @active-panel]))
