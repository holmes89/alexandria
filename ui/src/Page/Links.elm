module Page.Links exposing (Model, Msg, init, update, view)

import Html exposing (..)
import Html.Attributes exposing (class, href, id, placeholder, rows, src, style, target, value)
import Html.Events exposing (onClick, onInput)
import Http
import Links exposing (Link, linkDecoder, linkEncoder, linksDecoder)
import Session exposing (..)
import Time exposing (Posix)



-- MODEL


type Status
    = Failure
    | Loading
    | Success


type alias Model =
    { token : Token
    , status : Status
    , list : List Link
    , content : String
    }


init : Token -> ( Model, Cmd Msg )
init token =
    ( { token = token
      , status = Loading
      , list = []
      , content = ""
      }
    , fetchLinks token
    )



-- UPDATE


type Msg
    = FetchLinks (Result Http.Error (List Link))
    | AddLink (Result Http.Error Link)
    | Error
    | UpdateContent String
    | SendContent


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        FetchLinks result ->
            case result of
                Ok list ->
                    ( { model | status = Success, list = list }, Cmd.none )

                Err _ ->
                    ( { model | status = Failure }, Cmd.none )

        AddLink result ->
            case result of
                Ok entry ->
                    ( { model | status = Success, list = entry :: model.list, content = "" }, Cmd.none )

                Err _ ->
                    ( { model | status = Failure }, Cmd.none )

        Error ->
            ( { model | status = Failure }, Cmd.none )

        UpdateContent content ->
            ( { model | content = content }, Cmd.none )

        SendContent ->
            ( model, createLink model )



-- SUBSCRIPTIONS


subscriptions : Model -> Sub Msg
subscriptions model =
    Sub.none


viewLink : Link -> Html Msg
viewLink link =
    div [ class "column is-full" ]
        [ div [ class "card" ]
            [ div [ class "card-content" ]
                [ article [ class "media" ]
                    [ figure [ class "media-left " ]
                        [ p [ class "image is-64x64" ]
                            [ img [ src link.iconPath ] [] ]
                        ]
                    , div [ class "media-content" ]
                        [ div [ class "content link" ]
                            [ a [ href link.link, target "_blank" ] [ h4 [] [ text link.displayName ] ] ]
                        ]
                    ]
                ]
            ]
        ]


view : Model -> Html Msg
view model =
    case ( model.status, model.list ) of
        ( Failure, _ ) ->
            div []
                [ text "Failed"
                ]

        ( Loading, _ ) ->
            text "Loading..."

        ( Success, entries ) ->
            section [ class "section" ]
                [ div [ class "container" ]
                    [ div [ class "columns is-centered" ]
                        [ div [ class "column is-full" ]
                            [ div [ class "field has-addons has-text-centered" ]
                                [ div [ class "control has-icons-left  add-link" ]
                                    [ input [ class "input", value model.content, placeholder "Link", onInput UpdateContent ] []
                                    , span [ class "icon is-small is-left" ]
                                        [ i [ class "fas fa-link" ] [] ]
                                    ]
                                , div [ class "control" ]
                                    [ button [ class "button is-dark", onClick SendContent ] [ text "Submit" ]
                                    ]
                                ]
                            ]
                        ]
                    , div [ class "columns is-multiline" ]
                        (List.map
                            viewLink
                            entries
                        )
                    ]
                ]



-- HTTP


fetchLinks : Token -> Cmd Msg
fetchLinks token =
    Http.request
        { body = Http.emptyBody
        , expect = Http.expectJson FetchLinks linksDecoder
        , headers = [ Http.header "Authorization" token ]
        , method = "GET"
        , timeout = Nothing
        , tracker = Nothing
        , url = "https://docs.jholmestech.com/links/"
        }


createLink : Model -> Cmd Msg
createLink model =
    let
        link =
            { id = "", link = model.content, displayName = "", iconPath = "", created = Time.millisToPosix 0 }

        token =
            model.token
    in
    Http.request
        { body = Http.jsonBody <| linkEncoder link
        , expect = Http.expectJson AddLink linkDecoder
        , headers = [ Http.header "Authorization" token ]
        , method = "POST"
        , timeout = Nothing
        , tracker = Nothing
        , url = "https://docs.jholmestech.com/links/"
        }
