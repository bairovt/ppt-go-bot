package main

import (
	"regexp"
	"strings"

	arangoDriver "github.com/arangodb/go-driver"
)

func checkRegexPoints(body string) ( newBody string, err error) {
  // todo: optimize this
	var cursor arangoDriver.Cursor
	var point PointDoc	

  newBody = body
	queryPointsWithRegex := `
		FOR point IN Points
			FILTER TO_BOOL(point.regex)
			RETURN point
		`
  cursor, err = adb.Query(nil, queryPointsWithRegex, nil);
	if err != nil {
		return "", err
	}
	defer cursor.Close()
	// regexPoints
  // for (let point of regexPoints) {
  for {
		_, err = cursor.ReadDocument(nil, &point)

		if arangoDriver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return "", err
		} 
		
		re := regexp.MustCompile(point.Regex)
		re.ReplaceAllString(newBody, " " + point.Names[0] + " ")
    // newBody = newBody.replace(new RegExp(point.regex, 'g'), ' ' + point.names[0] + ' ');
  }
  return newBody, nil
}

func cleanBodyParser(body string) (str string, err error) {
	str = " " + strings.ToLower(body) + " "
	re := regexp.MustCompile(`‐|—`)
	str = re.ReplaceAllString(str, "-")
	

  // str = str.replace(/‐|—/g, '-');

  regexList := []*regexp.Regexp{
    regexp.MustCompile(`[\.,\?？:!\\\/\(\)\+&_"]|\w|\d|\n`),
    regexp.MustCompile(`здрав\S*|пр(и|е)ве(т|д)\S*|добр(ый|ое)|ч(е|и)л(о|а)век\S*|\sчел\s`),
    regexp.MustCompile(`\sищу\s|машин\S|\возь?м\S*|пас+аж\S*|меня|\s(по|вы)?ед(у|ет|ит|ем|им)|выез\S*|(у|до|по|за)е(ду|хать|дем)`),
    regexp.MustCompile(`могу|взять|(не)?больш(ой|ие|ые)|посылк\S*|сумк\S*|туда|(о|а)братно|чере(з|с)|дан+ый|момент\S*|срочно`),
    regexp.MustCompile(`\sс(ей|е|и|ий)час|сегод\S*|после|обед\S*|завтра|утр\S*|дн(е|ё)м|вечер(ом)?|дн(и|я)|день|\S*ночь\S*`),
    regexp.MustCompile(`\sв?т(е|и)чен\S*|час\S*|\sмин\s|минут\S*|пр(и|е)мер\S*|где(\s*-*\s*)то|(о|а)р(ие|и|е)нт(и|е)р\S*`),
    regexp.MustCompile(`\sкто\s|(кто|кем)?\-?\s*нибудь|мож(е|и)т|ближ\S*|врем\S*|увeз\S*|в?пут(ь|и)|п(о|а)пут\S*|есть|мест\S*`),
    regexp.MustCompile(`\sод(но|ин)\S*|дв(а|ое|ух|оих)|тр(и|ое)|номер\S*|\sтел\S*\s|\sзвонит.\s|прям\S*`),
    regexp.MustCompile(`декбр\S*|январ\S*|март\S*|апрел\S*|мая\s|июн\S*|июл\S*|август\S*|сентябр\S*|ноябрь\S*`),
    regexp.MustCompile(`\sпн\s|\sвт\s|\sср\s|\sчт\s|\sпт\s|\sсб\s|\sвс\s|ч(и|е)сл\S*`),
    regexp.MustCompile(`п(о|а)н(е|и)д\S*|вторн\S*|сред\S*|ч(е|и)тв\S*|пятн\S*|суб+от\S|в(о|а)скр\S*`),
    regexp.MustCompile(`личк.|\sлс\s|реб(е|ё)н\S*|дет(и|ей)|цена|\S*плат\S*|\sруб\S*|писать|п(и|е)ш(и|ы)те`),
    regexp.MustCompile(`нужн\S*|жела\S*|н(е|и)обх\S*|остал\S*|сторон\S*|заран(е|и)|(от|раз)дельн\S*`),
    regexp.MustCompile(`в(а|о)тсап\S*|вайбер\S*|телеграм\S*|тел(\s|$)|первой|второй|п(о|а)л(о|а)вин\S*`),
    regexp.MustCompile(`п(о|а)жалуй?ста\S*|легк\S*|авто\S*|багаж\S*|груз\S*|водител\S*|кресл\S`),
    regexp.MustCompile(`можно|чуть|(по)?раньше`),
    regexp.MustCompile(`тойот\S*|нис+ан\S*|хонд\S*|комфорт\S*`),
    regexp.MustCompile(`\sили\s|\sдля\s|\sбез\s`),
	}
  
	// let regexList = [
  //   /[\.,\?？:!\\\/\(\)\+&_"]|\w|\d|\n/g,
  //   /здрав\S*|пр(и|е)ве(т|д)\S*|добр(ый|ое)|ч(е|и)л(о|а)век\S*|\sчел\s/g,
  //   /\sищу\s|машин\S|\возь?м\S*|пас+аж\S*|меня|\s(по|вы)?ед(у|ет|ит|ем|им)|выез\S*|(у|до|по|за)е(ду|хать|дем)/g,
  //   /могу|взять|(не)?больш(ой|ие|ые)|посылк\S*|сумк\S*|туда|(о|а)братно|чере(з|с)|дан+ый|момент\S*|срочно/g,
  //   /\sс(ей|е|и|ий)час|сегод\S*|после|обед\S*|завтра|утр\S*|дн(е|ё)м|вечер(ом)?|дн(и|я)|день|\S*ночь\S*/g,
  //   /\sв?т(е|и)чен\S*|час\S*|\sмин\s|минут\S*|пр(и|е)мер\S*|где(\s*-*\s*)то|(о|а)р(ие|и|е)нт(и|е)р\S*/g,
  //   /\sкто\s|(кто|кем)?\-?\s*нибудь|мож(е|и)т|ближ\S*|врем\S*|увeз\S*|в?пут(ь|и)|п(о|а)пут\S*|есть|мест\S*/g,
  //   /\sод(но|ин)\S*|дв(а|ое|ух|оих)|тр(и|ое)|номер\S*|\sтел\S*\s|\sзвонит.\s|прям\S*/g,
  //   /декбр\S*|январ\S*|март\S*|апрел\S*|мая\s|июн\S*|июл\S*|август\S*|сентябр\S*|ноябрь\S*/g,
  //   /\sпн\s|\sвт\s|\sср\s|\sчт\s|\sпт\s|\sсб\s|\sвс\s|ч(и|е)сл\S*/g,
  //   /п(о|а)н(е|и)д\S*|вторн\S*|сред\S*|ч(е|и)тв\S*|пятн\S*|суб+от\S|в(о|а)скр\S*/g,
  //   /личк.|\sлс\s|реб(е|ё)н\S*|дет(и|ей)|цена|\S*плат\S*|\sруб\S*|писать|п(и|е)ш(и|ы)те/g,
  //   /нужн\S*|жела\S*|н(е|и)обх\S*|остал\S*|сторон\S*|заран(е|и)|(от|раз)дельн\S*/g,
  //   /в(а|о)тсап\S*|вайбер\S*|телеграм\S*|тел(\s|$)|первой|второй|п(о|а)л(о|а)вин\S*/g,
  //   /п(о|а)жалуй?ста\S*|легк\S*|авто\S*|багаж\S*|груз\S*|водител\S*|кресл\S/g,
  //   /можно|чуть|(по)?раньше/g,
  //   /тойот\S*|нис+ан\S*|хонд\S*|комфорт\S*/g,
  //   /\sили\s|\sдля\s|\sбез\s/g,
  // ];

  
	// for (let regex of regexList) {
  //   str = str.replace(regex, ' ');
  // }
	for _, re = range regexList {
    str = re.ReplaceAllString(str, " ");
  }    

	// убираем пробелы рядом с тире
	// str = str.replace(/\s*\-\s*/g, '-'); // убираем пробелы рядом с тире
	re = regexp.MustCompile(`\s*\-\s*`) 
	str = re.ReplaceAllString(str, "-")  
  
	str, err = checkRegexPoints(str)
	if err != nil {
		return "", err
	}
  // str = await checkRegexPoints(str);
  // str = str.replace(/\-/g, ' '); // все тире на пробелы, т.к. парные назв уже обработаны

  // str = str.replace(/\s\S{1,2}\s/g, ' '); // 1-2 буквы между пробелами
  // str = str.trim();
  // str = str.replace(/(^\S{1,2}\s)|(\s\S{1,2}$)/g, ''); // одиночные буквы между пробелами или тире
  // str = str.replace(/\s{2,}/g, ' ');

  return strings.TrimSpace(str), nil
}