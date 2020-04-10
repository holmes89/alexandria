module Page.Journal exposing (Model, Msg, init, update, view)

import Html exposing (..)
import Html.Attributes exposing (class, href, id, rows, src, style, value)
import Html.Events exposing (onClick, onInput)
import Http
import Journal exposing (Entry, entriesDecoder, entryDecoder, entryEncoder)
import Session exposing (..)
import Time exposing (Month, Posix, toDay, toHour, toMinute, toMonth, toYear)
import TimeZone exposing (america__new_york)



-- MODEL


type Status
    = Failure
    | Loading
    | Success


type alias Model =
    { token : Token
    , status : Status
    , list : List Entry
    , content : String
    }


init : Token -> ( Model, Cmd Msg )
init token =
    ( { token = token
      , status = Loading
      , list = []
      , content = ""
      }
    , fetchEntries token
    )



-- UPDATE


type Msg
    = FetchEntries (Result Http.Error (List Entry))
    | AddEntry (Result Http.Error Entry)
    | Error
    | UpdateContent String
    | SendContent


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        FetchEntries result ->
            case result of
                Ok list ->
                    ( { model | status = Success, list = list }, Cmd.none )

                Err _ ->
                    ( { model | status = Failure }, Cmd.none )

        AddEntry result ->
            case result of
                Ok entry ->
                    ( { model | status = Success, list = model.list ++ [ entry ], content = "" }, Cmd.none )

                Err _ ->
                    ( { model | status = Failure }, Cmd.none )

        Error ->
            ( { model | status = Failure }, Cmd.none )

        UpdateContent content ->
            ( { model | content = content }, Cmd.none )

        SendContent ->
            ( model, createEntry model )



-- SUBSCRIPTIONS


subscriptions : Model -> Sub Msg
subscriptions model =
    Sub.none



-- VIEW


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
            section [ class "hero" ]
                [ div [ class "hero-body journal-entry-area" ]
                    [ div [ class "container" ]
                        (List.map
                            viewEntry
                            entries
                        )
                    ]
                , div [ class "hero-foot journal-draft-area" ]
                    [ div [ class "container" ]
                        [ div [ class "columns" ]
                            [ div [ class "column is-11" ]
                                [ textarea [ class "textarea", value model.content, onInput UpdateContent, rows 5 ] [] ]
                            , div [ class "column journal-entry-button" ]
                                [ button [ class "button", onClick SendContent ] [ text "Submit" ] ]
                            ]
                        ]
                    ]
                ]


viewEntry : Entry -> Html Msg
viewEntry entry =
    div [ class "columns" ]
        [ div [ class "column is-1 has-text-right" ]
            [ div [] [ text (formatDate entry.created) ]
            , div [] [ text (formatTime entry.created) ]
            ]
        , div [ class "column" ]
            [ article [ class "message is-dark", id entry.id ]
                [ div [ class "message-body" ]
                    (List.map
                        (\l -> div [ class "journal-entry-margin-top" ] [ text l ])
                        (String.lines entry.content)
                    )
                ]
            ]
        ]


formatDate : Posix -> String
formatDate time =
    let
        zone =
            america__new_york ()
    in
    monthToString (toMonth zone time) ++ "/" ++ String.fromInt (toDay zone time) ++ "/" ++ String.fromInt (toYear zone time)


formatTime : Posix -> String
formatTime time =
    let
        zone =
            america__new_york ()
    in
    String.fromInt (toHour zone time) ++ ":" ++ minuteToString (toMinute zone time)


minuteToString : Int -> String
minuteToString minute =
    case minute of
        1 ->
            "01"

        2 ->
            "02"

        3 ->
            "03"

        4 ->
            "04"

        5 ->
            "05"

        6 ->
            "06"

        7 ->
            "07"

        8 ->
            "08"

        9 ->
            "09"

        _ ->
            String.fromInt minute


monthToString : Month -> String
monthToString month =
    case month of
        Time.Jan ->
            "1"

        Time.Feb ->
            "2"

        Time.Mar ->
            "3"

        Time.Apr ->
            "4"

        Time.May ->
            "5"

        Time.Jun ->
            "6"

        Time.Jul ->
            "7"

        Time.Aug ->
            "8"

        Time.Sep ->
            "9"

        Time.Oct ->
            "10"

        Time.Nov ->
            "11"

        Time.Dec ->
            "12"



-- HTTP


fetchEntries : Token -> Cmd Msg
fetchEntries token =
    Http.request
        { body = Http.emptyBody
        , expect = Http.expectJson FetchEntries entriesDecoder
        , headers = [ Http.header "Authorization" token ]
        , method = "GET"
        , timeout = Nothing
        , tracker = Nothing
        , url = "https://docs.jholmestech.com/journal/entry/"
        }


createEntry : Model -> Cmd Msg
createEntry model =
    let
        entry =
            { id = "", content = model.content, created = Time.millisToPosix 0 }

        token =
            model.token
    in
    Http.request
        { body = Http.jsonBody <| entryEncoder entry
        , expect = Http.expectJson AddEntry entryDecoder
        , headers = [ Http.header "Authorization" token ]
        , method = "POST"
        , timeout = Nothing
        , tracker = Nothing
        , url = "https://docs.jholmestech.com/journal/entry/"
        }
