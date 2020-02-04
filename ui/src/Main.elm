module Main exposing(..)

import Browser
import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (..)
import Http
import Json.Decode exposing (Decoder, field, string)

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
  = MorePlease
  | GotBook (Result Http.Error (List Book))


update : Msg -> Model -> (Model, Cmd Msg)
update msg model =
  case msg of
    MorePlease ->
      (Loading, getBook)

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
    [ h2 [] [ text "Books" ]
    , viewBook model
    ]


viewBook : Model -> Html Msg
viewBook model =
  case model of
    Failure ->
      div []
        [ text "Failed"
        ]

    Loading ->
      text "Loading..."

    Success docs ->
      div []
        (List.map(\l -> div [] [text l.displayName]) docs)



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
  Json.Decode.map2 Book (field "id" string) (field "display_name" string)

listBookDecoder : Decoder (List Book)
listBookDecoder =
  Json.Decode.list bookDecoder
