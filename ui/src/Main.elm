module Main exposing(..)

import Browser
import Html exposing (..)
import Html.Attributes exposing (class, style, src, href)
import Html.Events exposing (..)
import Http
import Json.Decode as Decode exposing (Decoder, field, string)
import Json.Decode.Pipeline exposing (required, optional)


-- MAIN


main =
  Browser.element
    { init = init
    , update = update
    , subscriptions = subscriptions
    , view = view
    }



-- MODEL


type Model
  = Failure
  | Loading
  | Success (List Book)


init : () -> (Model, Cmd Msg)
init _ =
  (Loading, getBook)



-- UPDATE


type Msg
  =  GotBook (Result Http.Error (List Book))


update : Msg -> Model -> (Model, Cmd Msg)
update msg model =
  case msg of

    GotBook result ->
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
            span [] [text "Alexandria", img [src "/assets/alexandria.png"] []]
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

    Success docs ->
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
            ]) docs)


-- HTTP


getBook : Cmd Msg
getBook =
  Http.get
    { url = "https://alexandria-api-4josd7vm2q-ue.a.run.app/documents/"
    , expect = Http.expectJson GotBook listBookDecoder
    }

type alias Book = {id : String, displayName : String}
bookDecoder : Decoder Book
bookDecoder =
  Decode.succeed Book
    |> required "id" string
    |> required "display_name" string

listBookDecoder : Decoder (List Book)
listBookDecoder =
  Decode.list bookDecoder
