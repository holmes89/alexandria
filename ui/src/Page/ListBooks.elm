module Page.ListBooks exposing (Model, Msg, init, update, view)

import Html exposing (..)
import Html.Attributes exposing (class, style, src, href)
import Html.Events exposing (..)
import Http
import Book exposing(Book, booksDecoder)

-- MODEL
type Model
  = Failure
  | Loading
  | Success (List Book)


init : () -> (Model, Cmd Msg)
init _ =
  (Loading, fetchBooks)

-- UPDATE
type Msg
  =  FetchBooks (Result Http.Error (List Book))

update : Msg -> Model -> (Model, Cmd Msg)
update msg model =
  case msg of
    FetchBooks result ->
      case result of
        Ok url ->
          (Success url, Cmd.none)
        Err _ ->
          (Failure, Cmd.none)

-- SUBSCRIPTIONS
subscriptions : Model -> Sub Msg
subscriptions model =
  Sub.none


-- VIEW
view : Model -> Html Msg
view model =
  div []
    [ nav [class "navbar", class "is-light"] [
      div [class "navbar-brand"] [
          div [class "navbar-item"] [
            span [] [text "Alexandria", img [src "/alexandria.png"] []]
          ]
        ]
      ]
    , div [] [
      section [class "section"] [
        div [class "container"] [
          viewBooks model
        ]
      ]
    ]
  ]

viewBooks : Model -> Html Msg
viewBooks model =
  case model of
    Failure ->
      div []
        [ text "Failed"
        ]

    Loading ->
      text "Loading..."

    Success books ->
        div [class "columns", class "is-mobile", class "is-multiline"]
            (List.map(\l -> div [class "column", class "is-one-quarter"] [
              div [class "card"] [
                header [class "card-header"] [
                  p [class "card-header-title"] [text l.displayName]
                ],
                div [class "card-content", style "text-align" "center"] [
                  img [src ("http://read.jholmestech.com/assets/covers/"++l.id++".jpg"), style "max-width" "300px"] []
                ],
                footer [class "card-footer"] [
                  a [class "card-footer-item", href "#"] [i [class "fas", class "fa-book-open"] [], text "Read"]
                ]
              ]
            ]) books)


-- HTTP


fetchBooks : Cmd Msg
fetchBooks =
  Http.get
    { url = "https://docs.jholmestech.com/books/"
    , expect = Http.expectJson FetchBooks booksDecoder
    }
